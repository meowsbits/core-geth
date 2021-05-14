// Copyright 2017 The go-ethereum Authors
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

package keccak

import (
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

type diffTest struct {
	ParentTimestamp    uint64
	ParentDifficulty   *big.Int
	CurrentTimestamp   uint64
	CurrentBlocknumber *big.Int
	CurrentDifficulty  *big.Int
}

func (d *diffTest) UnmarshalJSON(b []byte) (err error) {
	var ext struct {
		ParentTimestamp    string
		ParentDifficulty   string
		CurrentTimestamp   string
		CurrentBlocknumber string
		CurrentDifficulty  string
	}
	if err := json.Unmarshal(b, &ext); err != nil {
		return err
	}

	d.ParentTimestamp = math.MustParseUint64(ext.ParentTimestamp)
	d.ParentDifficulty = math.MustParseBig256(ext.ParentDifficulty)
	d.CurrentTimestamp = math.MustParseUint64(ext.CurrentTimestamp)
	d.CurrentBlocknumber = math.MustParseBig256(ext.CurrentBlocknumber)
	d.CurrentDifficulty = math.MustParseBig256(ext.CurrentDifficulty)

	return nil
}

func TestCalcDifficulty(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "..", "tests", "testdata", "BasicTests", "difficulty.json"))
	if err != nil {
		t.Skip(err)
	}
	defer file.Close()

	tests := make(map[string]diffTest)
	err = json.NewDecoder(file).Decode(&tests)
	if err != nil {
		t.Fatal(err)
	}

	config := &goethereum.ChainConfig{HomesteadBlock: big.NewInt(1150000)}

	for name, test := range tests {
		number := new(big.Int).Sub(test.CurrentBlocknumber, big.NewInt(1))
		diff := CalcDifficulty(config, test.CurrentTimestamp, &types.Header{
			Number:     number,
			Time:       test.ParentTimestamp,
			Difficulty: test.ParentDifficulty,
		})
		if diff.Cmp(test.CurrentDifficulty) != 0 {
			t.Error(name, "failed. Expected", test.CurrentDifficulty, "and calculated", diff)
		}
	}
}

// TestAstorDifficulty tests Astor testnet difficulty calculations.
// This references an issue reported here: https://github.com/etclabscore/core-geth/pull/369#issuecomment-835110086
/*
WARN [05-01|00:32:33.211] Invalid header encountered               number=1 hash="713364…241bb3" parent="f9c534…15f5bc" err="invalid difficulty: have 1046361407488, want 1098974756864"
DEBUG[05-01|00:32:33.211] Block body download terminated           err="syncing canceled (requested)"
DEBUG[05-01|00:32:33.211] Transaction receipt download terminated  err="syncing canceled (requested)"
*/
/*
For reference, this is genesis from Besu's eth_getBlockByNumber("earliest",false).
This is used to confirm the validity of the genesis block encoding.
> docker run -p 8545:8545 hyperledger/besu --rpc-http-enabled --rpc-http-api=admin,eth,debug,miner,net,txpool,priv,trace,web3 --rpc-http-cors-origins="all"

{
  "jsonrpc" : "2.0",
  "id" : 22236,
  "result" : {
    "number" : "0x0",
    "hash" : "0xf9c534ad514bc380e23a74e788c3d1f2446f9ad1cb74f0dbb48ea56c0315f5bc",
    "mixHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce" : "0x0000000000000042",
    "sha3Uncles" : "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "logsBloom" : "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "transactionsRoot" : "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "stateRoot" : "0x4160a1e20a109b2b0545add2299a5b425339f30e30becfa5afbaa9eae88c46cb",
    "receiptsRoot" : "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "miner" : "0x0000000000000000000000000000000000000000",
    "difficulty" : "0x10000000000",
    "totalDifficulty" : "0x10000000000",
    "extraData" : "0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa",
    "size" : "0x222",
    "gasLimit" : "0x1fffffffffffff",
    "gasUsed" : "0x0",
    "timestamp" : "0x0",
    "uncles" : [ ],
    "transactions" : [ ]
  }
}
*/

func TestAstorDifficulty(t *testing.T) {
	astorGenesisJSON := `{
  "config": {
    "chainId": 212,
    "ecip1049Block": 0,
    "keccak256": {
    }
  },
  "nonce": "0x42",
  "timestamp": "0x0",
  "extraData": "0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa",
  "gasLimit": "0x1fffffffffffff",
  "difficulty": "0x10000000000",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "6e3da3dd22f043958fdb862db5876e5e52a3d7b7": {
      "balance": "90000000000000000000000"
    },
    "46652437dc2f978952c5f5667533874c294e4e84": {
      "balance": "90000000000000000000000"
    },
    "3d650ae1134709a9f2f1c7785cd99400482cdd97": {
      "balance": "90000000000000000000000"
    },
    "e6b2d417af0e3686982dba572fcda575a4e8f575": {
      "balance": "90000000000000000000000"
    },
    "d118ab56776b158c70a0d03ed5ff03392a63deaf": {
      "balance": "90000000000000000000000"
    },
    "2bac997f3c4a51e655d14733aa1d75817443779a": {
      "balance": "90000000000000000000000"
    },
    "bd74228dca706ab4ab7d8b856ec6700624420831": {
      "balance": "90000000000000000000000"
    },
    "33bc5fb03008d7f0001638f48bfff2d628d16882": {
      "balance": "90000000000000000000000"
    },
    "ef826b2e0f6bf62edfe628dbeaa1af58f40f3b06": {
      "balance": "90000000000000000000000"
    },
    "7defc96deecc05288ff177ca08c025ed5c139a95": {
      "balance": "90000000000000000000000"
    },
    "2958DB51a0b4c458d0aa183E8cFB4f2E95cf6E75": {
      "balance": "90000000000000000000000"
    },
    "B460ce2C10d959251792D72b7E2B3EC684F013f1": {
      "balance": "90000000000000000000000"
    },
    "25848A733da7c782fB963D2e562cF7246bD5b6df": {
      "balance": "90000000000000000000000"
    }
  }
}`

	genesis := &genesisT.Genesis{Config: &coregeth.CoreGethChainConfig{}}
	err := json.Unmarshal([]byte(astorGenesisJSON), genesis)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	genesisBlock := core.GenesisToBlock(genesis, nil)

	// block1Time is an estimate of Astor testnet's block #1 time.
	// We can use `now` because the genesis block's time is 0; so the difficulty algo will have max'd out equivalently
	// for `actual block 1 timestamp` and `now`.
	block1Time := uint64(time.Now().Unix())

	// Values taken from reported logs (see test comment above).
	logCalculated := big.NewInt(1098974756864)
	astor := big.NewInt(1046361407488)

	got := CalcDifficulty(genesis.Config, block1Time, genesisBlock.Header())
	if got.Cmp(logCalculated) != 0 {
		t.Fatalf("log 'want' value mismatches test value; should not fail ever (unless logs or test is wrong)")
	}

	if got.Cmp(astor) != 0 {
		t.Log("confirm") // Our test value matches the logged 'want' value.
	}

	if genesisBlock.Hash() != common.HexToHash("0xf9c534ad514bc380e23a74e788c3d1f2446f9ad1cb74f0dbb48ea56c0315f5bc") {
		t.Fatalf("hash mismatch, want: %s", genesisBlock.Hash().Hex())
	}
}
