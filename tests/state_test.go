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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

func TestState(t *testing.T) {
	//t.Parallel()

	st := new(testMatcher)
	// Long tests:
	// st.whitelist(`^stAttackTest/ContractCreationSpam`)
	// st.whitelist(`^stBadOpcode/badOpcodes`)
	// st.whitelist(`^stPreCompiledContracts/modexp`)
	// st.whitelist(`^stQuadraticComplexityTest/`)
	// st.whitelist(`^stStaticCall/static_Call50000`)
	// st.whitelist(`^stStaticCall/static_Return50000`)
	// st.whitelist(`^stSystemOperationsTest/CallRecursiveBomb`)
	// st.whitelist(`^stTransactionTest/Opcodes_TransactionInit`)

	// Very time consuming
	// st.whitelist(`^stTimeConsuming/`)

	// Uses 1GB RAM per tested fork
	// st.whitelist(`^stStaticCall/static_Call1MB`)

	// Broken tests:
	// Expected failures:
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/3`, "bug in test")

	heads := make(chan *types.Header)
	//var sub ethereum.Subscription
	var err error
	quit := make(chan bool)

	// if client has met istanbul and we have run tests against it,
	// then we can exit.
	metIstanbul := false
	nextFork := make(chan struct{}, 1)

	forkEnabled := func(name string, num uint64) bool {
		na := ""
		for i, v := range fnames {
			n := fblocks[i]
			if num >= n {
				na = v
			}
		}
		return name == na
	}

	fnameForClient := func(n uint64) string {
		name := ""
		for i, v := range fblocks {
			if n >= v {
				name = fnames[i]
				continue
			}
			break
		}
		return name
	}

	if os.Getenv("AM") != "" {
		MyTransmitter = NewTransmitter()
		_, err = MyTransmitter.client.SubscribeNewHead(context.Background(), heads)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Transmitter OKGO")
		t.Log("log: Transmitter OKGO")

		//defer sub.Unsubscribe()
	}

	if os.Getenv("AM") != "" {
		go func() {
			zeroTxsN := 0
			// didBump := false
			lfn := ""
			for {
				select {
				// case err := <-sub.Err():
				//	if err != nil {
				//		t.Fatal("subscription error", err)
				//	}
				case head := <-heads:
					bl, err := MyTransmitter.client.BlockByHash(MyTransmitter.ctx, head.Hash())
					if err != nil {
						t.Fatal(err)
					}
					fn := fnameForClient(bl.NumberU64())
					if lfn != fn {
						MyTransmitter.purgePending()
						lfn = fn
						nextFork <- struct{}{}
					}
					if bl.Transactions().Len() == 0 {
						zeroTxsN++
					} else {
						zeroTxsN = 0
					}
					if zeroTxsN > 10 && len(MyTransmitter.txPendingContracts) > 0 {
						// if !didBump {
						// 	to := common.HexToAddress("0xb1355e69c1ba7b401E38c95CdB936727aE88bE76")
						// 	n, _ := MyTransmitter.client.NonceAt(MyTransmitter.ctx, MyTransmitter.sender, nil)
						// 	bump := types.NewMessage(MyTransmitter.sender, &to, n, new(big.Int).SetUint64(13), vars.TxGasContractCreation, big.NewInt(vars.GWei), nil, false)
						// 	MyTransmitter.SendMessage(bump)
						// 	didBump = true
						// 	zeroTxsN = 0
						// } else {
						MyTransmitter.purgePending()
						// }
					}
					MyTransmitter.currentBlock = bl.NumberU64()
					fmt.Println("New block head", "num", bl.NumberU64(), "txlen", bl.Transactions().Len(), "hash", bl.Hash().Hex(), "q", len(MyTransmitter.txPendingContracts))
					for _, tr := range bl.Transactions() {

						MyTransmitter.mu.Lock()
						v, ok := MyTransmitter.txPendingContracts[tr.Hash()]
						MyTransmitter.mu.Unlock()
						if !ok {
							continue
						}

						fmt.Println("Matched pending transaction", tr.Hash().Hex())

						receipt, err := MyTransmitter.client.TransactionReceipt(MyTransmitter.ctx, tr.Hash())
						if err != nil {
							panic(fmt.Sprintf("receipt err=%v tx=%x", err, tr.Hash()))
						}

						next := types.NewMessage(MyTransmitter.sender, &receipt.ContractAddress, 0, v.Value(), v.Gas(), v.GasPrice(), v.Data(), true)

						sentTxHash, err := MyTransmitter.SendMessage(next)
						if err != nil {
							fmt.Printf("FAILED to send pending tx, err=%v", err)
							// b, _ := json.MarshalIndent(next, "", "    ")
							// panic(fmt.Sprintf("send pend err=%v tx=%s", err, string(b)))
						}

						MyTransmitter.mu.Lock()
						delete(MyTransmitter.txPendingContracts, tr.Hash())
						MyTransmitter.mu.Unlock()

						fmt.Println("Sent was-pending tx", "hash", sentTxHash.Hex())
					}

				case <-quit:
					return

				}
			}
		}()
	}

	// For Istanbul, older tests were moved into LegacyTests
	lastForkName := ""
	for !metIstanbul {

		fname := fnameForClient(MyTransmitter.currentBlock)
		if lastForkName == fname {
			fmt.Printf("Finished dirs, same fork (%s) (waiting)\n", fname)
			<-nextFork
		}
		lastForkName = fnameForClient(MyTransmitter.currentBlock)

		dirs := []string{legacyStateTestDir}
		if fnameForClient(MyTransmitter.currentBlock) == "Istanbul" {
			dirs = append([]string{stateTestDir}, legacyStateTestDir)
		}
		for id, dir := range dirs {
			fmt.Printf("Running st walk: %d (client=%d => %s)", id, MyTransmitter.currentBlock, fnameForClient(MyTransmitter.currentBlock))
			st.walk(t, dir, func(t *testing.T, name string, test *StateTest) {
				for _, subtest := range test.Subtests() {
					subtest := subtest

					if !forkEnabled(subtest.Fork, MyTransmitter.currentBlock) {
						t.Logf("Skipping fork (%s > %d)", subtest.Fork, MyTransmitter.currentBlock)
						return
					}
					if subtest.Fork == "Istanbul" {
						metIstanbul = true
					}

					key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
					name := name + "/" + key

					fmt.Println("running test", name)
					t.Run(key+"/trie", func(t *testing.T) {
						withFakeTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
							_, _, err := test.Run(subtest, vmconfig, false)
							return st.checkFailure(t, name+"/trie", err)
						})
					})
					// time.Sleep(time.Second)
					// t.Run(key+"/snap", func(t *testing.T) {
					//	withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
					//		snaps, statedb, err := test.Run(subtest, vmconfig, true)
					//		if _, err := snaps.Journal(statedb.IntermediateRoot(false)); err != nil {
					//			return err
					//		}
					//		return st.checkFailure(t, name+"/snap", err)
					//	})
					// })
				}
			})
			MyTransmitter.wg.Wait()
			// t.Run("wait for pending txs", func(t *testing.T) {
			//	time.Sleep(5*time.Second)
			//	for len(MyTransmitter.txPendingContracts) > 0 {
			//		fmt.Println("Sleeping")
			//		time.Sleep(time.Minute)
			//		//MyTransmitter.mu.Lock()
			//		//for k, v := range MyTransmitter.txPendingContracts {
			//		//
			//		//}
			//		//
			//		//MyTransmitter.mu.Unlock()
			//		//
			//	}
			//	quit <- true
			//	close(quit)
			// })
		}
	}
	close(nextFork)
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withFakeTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	// Use config from command line arguments.
	config := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}
	test(config)
}

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	// Use config from command line arguments.
	config := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}
	err := test(config)
	if err == nil {
		return
	}

	// Test failed, re-run with tracing enabled.
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	tracer := vm.NewJSONLogger(&vm.LogConfig{DisableMemory: true}, w)
	config.Debug, config.Tracer = true, tracer
	err2 := test(config)
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	w.Flush()
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
	//t.Logf("EVM output: 0x%x", tracer.Output())
	//t.Logf("EVM error: %v", tracer.Error())
}
