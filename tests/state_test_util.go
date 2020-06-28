// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// StateTest checks transaction processing without block context.
// See https://github.com/ethereum/EIPs/issues/176 for the test format specification.
type StateTest struct {
	// Name is only used for writing tests (Marshaling).
	// It is set by inference, using the test's filepath base.
	Name string
	json stJSON
}

// StateSubtest selects a specific configuration of a General State Test.
type StateSubtest struct {
	Fork  string
	Index int
}

func (t *StateTest) UnmarshalJSON(in []byte) error {
	return json.Unmarshal(in, &t.json)
}

func (t StateTest) MarshalJSON() ([]byte, error) {
	m := map[string]stJSON{
		t.Name: t.json,
	}
	return json.MarshalIndent(m, "", "    ")
}

type stJSON struct {
	Info stInfo                   `json:"_info"`
	Env  stEnv                    `json:"env"`
	Pre  genesisT.GenesisAlloc    `json:"pre"`
	Tx   stTransaction            `json:"transaction"`
	Out  hexutil.Bytes            `json:"out"`
	Post map[string][]stPostState `json:"post"`
}

type stInfo struct {
	Comment            string         `json:"comment"`
	FillingRPCServer   string         `json:"filling-rpc-server"`
	FillingToolVersion string         `json:"filling-tool-version"`
	FilledWith         string         `json:"filledWith"`
	LLLCVersion        string         `json:"lllcversion"`
	Source             string         `json:"source"`
	SourceHash         string         `json:"sourceHash"`
	WrittenWith        string         `json:"written-with"`
	Parent             string         `json:"parent"`
	ParentSha1Sum      string         `json:"parentSha1Sum"`
	Chainspecs         chainspecRefsT `json:"chainspecs"`
}

type stPostState struct {
	Root    common.UnprefixedHash `json:"hash"`
	Logs    common.UnprefixedHash `json:"logs"`
	Indexes struct {
		Data  int `json:"data"`
		Gas   int `json:"gas"`
		Value int `json:"value"`
	}
}

//go:generate gencodec -type stEnv -field-override stEnvMarshaling -out gen_stenv.go

type stEnv struct {
	Coinbase   common.Address `json:"currentCoinbase"   gencodec:"required"`
	Difficulty *big.Int       `json:"currentDifficulty" gencodec:"required"`
	GasLimit   uint64         `json:"currentGasLimit"   gencodec:"required"`
	Number     uint64         `json:"currentNumber"     gencodec:"required"`
	Timestamp  uint64         `json:"currentTimestamp"  gencodec:"required"`
}

type stEnvMarshaling struct {
	Coinbase   common.UnprefixedAddress
	Difficulty *math.HexOrDecimal256
	GasLimit   math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
}

//go:generate gencodec -type stTransaction -field-override stTransactionMarshaling -out gen_sttransaction.go

type stTransaction struct {
	GasPrice   *big.Int `json:"gasPrice"`
	Nonce      uint64   `json:"nonce"`
	To         string   `json:"to"`
	Data       []string `json:"data"`
	GasLimit   []uint64 `json:"gasLimit"`
	Value      []string `json:"value"`
	PrivateKey []byte   `json:"secretKey"`
}

type stTransactionMarshaling struct {
	GasPrice   *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	GasLimit   []math.HexOrDecimal64
	PrivateKey hexutil.Bytes
}

// getVMConfig takes a fork definition and returns a chain config.
// The fork definition can be
// - a plain forkname, e.g. `Byzantium`,
// - a fork basename, and a list of EIPs to enable; e.g. `Byzantium+1884+1283`.
func getVMConfig(forkString string) (baseConfig ctypes.ChainConfigurator, eips []int, err error) {
	var (
		splitForks            = strings.Split(forkString, "+")
		ok                    bool
		baseName, eipsStrings = splitForks[0], splitForks[1:]
	)
	if baseConfig, ok = Forks[baseName]; !ok {
		return nil, nil, UnsupportedForkError{baseName}
	}
	for _, eip := range eipsStrings {
		if eipNum, err := strconv.Atoi(eip); err != nil {
			return nil, nil, fmt.Errorf("syntax error, invalid eip number %v", eipNum)
		} else {
			eips = append(eips, eipNum)
		}
	}
	return baseConfig, eips, nil
}

// Subtests returns all valid subtests of the test.
func (t *StateTest) Subtests() []StateSubtest {
	var sub []StateSubtest
	for fork, pss := range t.json.Post {
		for i := range pss {
			sub = append(sub, StateSubtest{fork, i})
		}
	}
	return sub
}

type Transmitter struct {
	ctx                context.Context
	client             *ethclient.Client
	sender             common.Address
	nonce              uint64
	currentBlock       uint64
	chainId            *big.Int
	ks                 *keystore.KeyStore
	txPendingContracts map[common.Hash]core.Message
	mu                 sync.Mutex
	wg                 sync.WaitGroup
	count              int
}

var MyTransmitter *Transmitter

func setupClient(ctx context.Context) (*ethclient.Client, *big.Int) {
	endpoint := os.Getenv("AM_CLIENT")
	if endpoint == "" {
		fmt.Println("NO ENDPOINT")
		os.Exit(1)
	}
	cl, err := ethclient.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	cid, err := cl.ChainID(ctx)
	if err != nil {
		panic(err)
	}
	return cl, cid
}

func setupKeystore() *keystore.KeyStore {
	keydir := os.Getenv("AM_KEYSTORE")
	if keydir == "" {
		panic("NO KEYSTORE")
	}
	ks := keystore.NewKeyStore(keydir, keystore.StandardScryptN, keystore.StandardScryptP)

	acc := os.Getenv("AM_ADDRESS")
	if acc == "" {
		panic("NO SENDING ADDRESS")
	}
	passFile := os.Getenv("AM_PASSWORDFILE")
	if passFile == "" {
		panic("NO PASSFILE")
	}
	password, err := ioutil.ReadFile(passFile)
	if err != nil {
		panic(err)
	}
	sender := common.HexToAddress(acc)
	err = ks.Unlock(accounts.Account{Address: sender}, strings.TrimSuffix(string(password), "\n"))
	if err != nil {
		panic(err)
	}
	return ks
}
func NewTransmitter() *Transmitter {
	ctx := context.Background()
	cl, cid := setupClient(ctx)
	ks := setupKeystore()
	sender := ks.Accounts()[0]

	bal, _ := cl.BalanceAt(ctx, sender.Address, nil)
	log.Printf("sender=%s balance=%d", sender.Address.String(), bal)

	t := &Transmitter{
		ctx:                ctx,
		client:             cl,
		currentBlock:       0,
		chainId:            cid,
		ks:                 ks,
		sender:             sender.Address,
		txPendingContracts: make(map[common.Hash]core.Message),
	}
	t.setNonce()
	return t
}

func (t *Transmitter) setNonce() {
	t.mu.Lock()
	defer t.mu.Unlock()
	nonce, err := t.client.NonceAt(t.ctx, t.sender, nil)
	if err != nil {
		log.Fatal("pending nonce", err)
	}
	t.nonce = uint64(nonce)
	fmt.Println("set transmitter nonce", nonce)
}

var maxValue = big.NewInt(vars.GWei) // As much as is willing to go into value in transaction.

var eip2b = uint64(115)
var eip155b = uint64(267)
var eip2028b = uint64(906)

func (t *Transmitter) SendMessage(msg core.Message) (common.Hash, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	defer t.wg.Done()

	txhash := common.Hash{}
	var tx *types.Transaction

	gas := msg.Gas()


	if t.currentBlock == 0 {
		b, _ := t.client.BlockByNumber(t.ctx, nil)
		t.currentBlock = b.NumberU64()
	}
	isEIP2, isEIP155, isEIP2028 := t.currentBlock >= eip2b, t.currentBlock >= eip155b, t.currentBlock >= eip2028b

	igas, err := core.IntrinsicGas(msg.Data(), msg.To() == nil, isEIP2, isEIP2028)
	if err != nil {
		panic(fmt.Sprintf("instrinsict gas calc err=%v", err))
	}
	gas += igas * 2
	if gas > 8000000 {
		gas = 8000000
	}

	val := msg.Value()
	if val.Cmp(maxValue) > 0 {
		val.Set(maxValue)
	}

	if msg.To() == nil {
		tx = types.NewContractCreation(t.nonce, val, gas, msg.GasPrice(), msg.Data())
	} else {
		to := *msg.To()
		tx = types.NewTransaction(t.nonce, to, val, gas, msg.GasPrice(), msg.Data())
		if strings.Contains(to.Hex(), "000000000000000000000000000000000") {
			b, _ := json.MarshalIndent(tx, "", "    ")
			panic(fmt.Sprintf("empty To address=%v", string(b)))
		}
	}

	var signed *types.Transaction
	if isEIP155 {
		signed, err = t.ks.SignTx(accounts.Account{Address: t.sender}, tx, t.chainId)
	} else {
		signed, err = t.ks.SignTx(accounts.Account{Address: t.sender}, tx, nil)
	}
	if err != nil {
		return txhash, fmt.Errorf("sign tx: %v", err)
	}
	err = t.client.SendTransaction(t.ctx, signed)

	//if t.count > 15 {
	//	time.Sleep(15 * time.Second)
	//	t.count = 0
	//}
	//t.count++
	//
	if err != nil {
		return txhash, fmt.Errorf("send tx: %v nonce=%d gas=%d gasPrice=%d value=%d", err, t.nonce, msg.Gas(), msg.GasPrice(), val)
	}
	t.nonce++
	return signed.Hash(), nil
}

// Run executes a specific subtest and verifies the post-state and logs
func (t *StateTest) Run(subtest StateSubtest, vmconfig vm.Config, snapshotter bool) (*snapshot.Tree, *state.StateDB, error) {
	snaps, statedb, root, err := t.RunNoVerify(subtest, vmconfig, snapshotter)
	if err != nil {
		return snaps, statedb, err
	}
	post := t.json.Post[subtest.Fork][subtest.Index]
	// N.B: We need to do this in a two-step process, because the first Commit takes care
	// of suicides, and we need to touch the coinbase _after_ it has potentially suicided.
	if root != common.Hash(post.Root) {
		return snaps, statedb, fmt.Errorf("post state root mismatch: got %x, want %x", root, post.Root)
	}
	if logs := rlpHash(statedb.Logs()); logs != common.Hash(post.Logs) {
		return snaps, statedb, fmt.Errorf("post state logs hash mismatch: got %x, want %x", logs, post.Logs)
	}
	return snaps, statedb, nil
}

// RunNoVerify runs a specific subtest and returns the statedb and post-state root
func (t *StateTest) RunNoVerify(subtest StateSubtest, vmconfig vm.Config, snapshotter bool) (*snapshot.Tree, *state.StateDB, common.Hash, error) {

	config, eips, err := getVMConfig(subtest.Fork)
	if err != nil {
		return nil, nil, common.Hash{}, UnsupportedForkError{subtest.Fork}
	}
	vmconfig.ExtraEips = eips
	block := core.GenesisToBlock(t.genesis(config), nil)
	snaps, statedb := MakePreState(rawdb.NewMemoryDatabase(), t.json.Pre, snapshotter)

	post := t.json.Post[subtest.Fork][subtest.Index]
	msg, err := t.json.Tx.toMessage(post)
	if err != nil {
		return nil, nil, common.Hash{}, err
	}
	context := core.NewEVMContext(msg, block.Header(), nil, &t.json.Env.Coinbase)
	context.GetHash = vmTestBlockHash
	evm := vm.NewEVM(context, statedb, config, vmconfig)

	gaspool := new(core.GasPool)
	gaspool.AddGas(block.GasLimit())
	snapshot := statedb.Snapshot()
	if _, err := core.ApplyMessage(evm, msg, gaspool); err != nil {
		statedb.RevertToSnapshot(snapshot)
	}

	n, err := MyTransmitter.client.TxpoolPending(MyTransmitter.ctx)

	//n, err := MyTransmitter.client.PendingTransactionCount(MyTransmitter.ctx)
	for ; err == nil && n > 1000; n, err = MyTransmitter.client.TxpoolPending(MyTransmitter.ctx) {
		fmt.Println("Txpool pending len=", n, "(sleeping)")
		time.Sleep(time.Second)
	}
	for ; err == nil && n > 1000; n, err = MyTransmitter.client.TxpoolQueued(MyTransmitter.ctx) {
		fmt.Println("Txpool queued len=", n, "(sleeping)")
		time.Sleep(time.Second)
	}

	//toHex := ""
	//if msg.To() != nil {
	//	toHex = msg.To().Hex()
	//}
	//fmt.Println("apply message", "msg", "from", msg.From().Hex(), "to", toHex, "data", common.Bytes2Hex(msg.Data()))

	if msg.From() == (common.Address{}) {
		panic("empty sender") // I don't think this will ever happen for the >=Istanbul test set.
	}

	// We need to deploy the target contract,
	// wait for it to be included in a block so we can get the transaction receipt (this is the 'pre' state),
	// then send the test transaction to the transactionReceipt.contractAddress.

	for preAccount, st := range t.json.Pre {
		// This is the constant sender address;
		/*
			> geth --keystore ./ks account import 0x45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8
			> geth --keystore ./ks account list
			INFO [05-29|11:53:58.909] Maximum peer count                       ETH=50 LES=0 total=50
			INFO [05-29|11:53:58.909] Smartcard socket not found, disabling    err="stat /run/pcscd/pcscd.comm: no such file or directory"
			Account #0: {a94f5374fce5edbc8e2a8697c15331677e6ebf0b} keystore:///home/ia/go/src/github.com/ethereum/go-ethereum/ks/UTC--2020-05-29T15-59-34.926438299Z--a94f5374fce5edbc8e2a8697c15331677e6ebf0b

		*/
		if preAccount == common.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b") {
			continue
		}
		if msg.To() == nil {
			continue
		}
		if *msg.To() != preAccount {
			continue
		}
		if len(st.Code) == 0 {
			continue
		}

		if MyTransmitter.currentBlock == 0 {
			b, _ := MyTransmitter.client.BlockByNumber(MyTransmitter.ctx, nil)
			MyTransmitter.currentBlock = b.NumberU64()
		}
		isEIP2, _, isEIP2028 := MyTransmitter.currentBlock >= eip2b, MyTransmitter.currentBlock >= eip155b, MyTransmitter.currentBlock >= eip2028b

		gas := uint64(0)
		igas, err := core.IntrinsicGas(msg.Data(), true, isEIP2, isEIP2028)
		if err != nil {
			panic(fmt.Sprintf("instrinsict gas calc err=%v", err))
		}
		gas += igas * 2
		if gas > 8000000 {
			gas = 8000000
		}

		tx := types.NewMessage(MyTransmitter.sender, nil, 0, new(big.Int), gas, big.NewInt(vars.GWei), st.Code, false)

		fmt.Println("PRE", t.Name, subtest.Fork, subtest.Index, preAccount.Hex(), common.Bytes2Hex(tx.Data()))
		MyTransmitter.wg.Add(1)
		txHash, err := MyTransmitter.SendMessage(tx)
		if err != nil {
			log.Fatal("send prestate message ", err)
		}

		// regsiter this sent transaction hash to the transaction that we should send when the one we just sent completes.
		MyTransmitter.mu.Lock()
		MyTransmitter.wg.Add(1)
		MyTransmitter.txPendingContracts[txHash] = msg
		MyTransmitter.mu.Unlock()

	}

	if msg.To() == nil {
		if MyTransmitter.currentBlock == 0 {
			b, _ := MyTransmitter.client.BlockByNumber(MyTransmitter.ctx, nil)
			MyTransmitter.currentBlock = b.NumberU64()
		}
		isEIP2, _, isEIP2028 := MyTransmitter.currentBlock >= eip2b, MyTransmitter.currentBlock >= eip155b, MyTransmitter.currentBlock >= eip2028b

		gas := uint64(0)
		igas, err := core.IntrinsicGas(msg.Data(), true, isEIP2, isEIP2028)
		if err != nil {
			panic(fmt.Sprintf("instrinsict gas calc err=%v", err))
		}
		gas += igas *2
		if gas > 8000000 {
			gas = 8000000
		}

		fmt.Println("TX", t.Name, subtest.Fork, subtest.Index, common.Bytes2Hex(msg.Data()))
		tx := types.NewMessage(MyTransmitter.sender, nil, 0, new(big.Int), gas, big.NewInt(vars.GWei), msg.Data(), false)
		MyTransmitter.wg.Add(1)
		txHash, err := MyTransmitter.SendMessage(tx)
		if err != nil {
			log.Fatal("send contract create ", err)
		}
		fmt.Println("sent contract creation", txHash.Hex())
	}

	// Commit block
	statedb.Commit(config.IsEnabled(config.GetEIP161dTransition, block.Number()))
	// Add 0-value mining reward. This only makes a difference in the cases
	// where
	// - the coinbase suicided, or
	// - there are only 'bad' transactions, which aren't executed. In those cases,
	//   the coinbase gets no txfee, so isn't created, and thus needs to be touched
	statedb.AddBalance(block.Coinbase(), new(big.Int))
	// And _now_ get the state root
	root := statedb.IntermediateRoot(config.IsEnabled(config.GetEIP161dTransition, block.Number()))
	return snaps, statedb, root, nil
}

func (t *StateTest) gasLimit(subtest StateSubtest) uint64 {
	return t.json.Tx.GasLimit[t.json.Post[subtest.Fork][subtest.Index].Indexes.Gas]
}

func MakePreState(db ethdb.Database, accounts genesisT.GenesisAlloc, snapshotter bool) (*snapshot.Tree, *state.StateDB) {
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb, nil)
	for addr, a := range accounts {
		statedb.SetCode(addr, a.Code)
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
		for k, v := range a.Storage {
			statedb.SetState(addr, k, v)
		}
	}
	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false)

	var snaps *snapshot.Tree
	if snapshotter {
		snaps = snapshot.New(db, sdb.TrieDB(), 1, root, false)
	}
	statedb, _ = state.New(root, sdb, snaps)
	return snaps, statedb
}

func (t *StateTest) genesis(config ctypes.ChainConfigurator) *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     config,
		Coinbase:   t.json.Env.Coinbase,
		Difficulty: t.json.Env.Difficulty,
		GasLimit:   t.json.Env.GasLimit,
		Number:     t.json.Env.Number,
		Timestamp:  t.json.Env.Timestamp,
		Alloc:      t.json.Pre,
	}
}

func (tx *stTransaction) toMessage(ps stPostState) (core.Message, error) {
	// Derive sender from private key if present.
	var from common.Address
	if len(tx.PrivateKey) > 0 {
		key, err := crypto.ToECDSA(tx.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %v", err)
		}
		from = crypto.PubkeyToAddress(key.PublicKey)
	}
	// Parse recipient if present.
	var to *common.Address
	if tx.To != "" {
		to = new(common.Address)
		if err := to.UnmarshalText([]byte(tx.To)); err != nil {
			return nil, fmt.Errorf("invalid to address: %v", err)
		}
	}

	// Get values specific to this post state.
	if ps.Indexes.Data > len(tx.Data) {
		return nil, fmt.Errorf("tx data index %d out of bounds", ps.Indexes.Data)
	}
	if ps.Indexes.Value > len(tx.Value) {
		return nil, fmt.Errorf("tx value index %d out of bounds", ps.Indexes.Value)
	}
	if ps.Indexes.Gas > len(tx.GasLimit) {
		return nil, fmt.Errorf("tx gas limit index %d out of bounds", ps.Indexes.Gas)
	}
	dataHex := tx.Data[ps.Indexes.Data]
	valueHex := tx.Value[ps.Indexes.Value]
	gasLimit := tx.GasLimit[ps.Indexes.Gas]
	// Value, Data hex encoding is messy: https://github.com/ethereum/tests/issues/203
	value := new(big.Int)
	if valueHex != "0x" {
		v, ok := math.ParseBig256(valueHex)
		if !ok {
			return nil, fmt.Errorf("invalid tx value %q", valueHex)
		}
		value = v
	}
	data, err := hex.DecodeString(strings.TrimPrefix(dataHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid tx data %q", dataHex)
	}

	msg := types.NewMessage(from, to, tx.Nonce, value, gasLimit, tx.GasPrice, data, true)
	return msg, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
