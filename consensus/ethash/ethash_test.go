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

package ethash

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// Tests that ethash works correctly in test mode.
func TestTestMode(t *testing.T) {
	header := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(100)}

	ethash := NewTester(nil, false)
	defer ethash.Close()

	results := make(chan *types.Block)
	err := ethash.Seal(nil, types.NewBlockWithHeader(header), results, nil)
	if err != nil {
		t.Fatalf("failed to seal block: %v", err)
	}
	select {
	case block := <-results:
		header.Nonce = types.EncodeNonce(block.Nonce())
		header.MixDigest = block.MixDigest()
		if err := ethash.VerifySeal(nil, header); err != nil {
			t.Fatalf("unexpected verification error: %v", err)
		}
	case <-time.NewTimer(2 * time.Second).C:
		t.Error("sealing result timeout")
	}
}

// This test checks that cache lru logic doesn't crash under load.
// It reproduces https://github.com/ethereum/go-ethereum/issues/14943
func TestCacheFileEvict(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "ethash-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	e := New(Config{CachesInMem: 3, CachesOnDisk: 10, CacheDir: tmpdir, PowMode: ModeTest}, nil, false)
	defer e.Close()

	workers := 8
	epochs := 100
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go verifyTest(&wg, e, i, epochs)
	}
	wg.Wait()
}

func verifyTest(wg *sync.WaitGroup, e *Ethash, workerIndex, epochs int) {
	defer wg.Done()

	const wiggle = 4 * epochLength
	r := rand.New(rand.NewSource(int64(workerIndex)))
	for epoch := 0; epoch < epochs; epoch++ {
		block := int64(epoch)*epochLength - wiggle/2 + r.Int63n(wiggle)
		if block < 0 {
			block = 0
		}
		header := &types.Header{Number: big.NewInt(block), Difficulty: big.NewInt(100)}
		e.VerifySeal(nil, header)
	}
}

func TestRemoteSealer(t *testing.T) {
	ethash := NewTester(nil, false)
	defer ethash.Close()

	api := &API{ethash}
	if _, err := api.GetWork(); err != errNoMiningWork {
		t.Error("expect to return an error indicate there is no mining work")
	}
	header := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(100)}
	block := types.NewBlockWithHeader(header)
	sealhash := ethash.SealHash(header)

	// Push new work.
	results := make(chan *types.Block)
	ethash.Seal(nil, block, results, nil)

	var (
		work [4]string
		err  error
	)
	if work, err = api.GetWork(); err != nil || work[0] != sealhash.Hex() {
		t.Error("expect to return a mining work has same hash")
	}

	if res := api.SubmitWork(types.BlockNonce{}, sealhash, common.Hash{}); res {
		t.Error("expect to return false when submit a fake solution")
	}
	// Push new block with same block number to replace the original one.
	header = &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(1000)}
	block = types.NewBlockWithHeader(header)
	sealhash = ethash.SealHash(header)
	ethash.Seal(nil, block, results, nil)

	if work, err = api.GetWork(); err != nil || work[0] != sealhash.Hex() {
		t.Error("expect to return the latest pushed work")
	}
}

func TestHashRate(t *testing.T) {
	var (
		hashrate = []hexutil.Uint64{100, 200, 300}
		expect   uint64
		ids      = []common.Hash{common.HexToHash("a"), common.HexToHash("b"), common.HexToHash("c")}
	)
	ethash := NewTester(nil, false)
	defer ethash.Close()

	if tot := ethash.Hashrate(); tot != 0 {
		t.Error("expect the result should be zero")
	}

	api := &API{ethash}
	for i := 0; i < len(hashrate); i += 1 {
		if res := api.SubmitHashRate(hashrate[i], ids[i]); !res {
			t.Error("remote miner submit hashrate failed")
		}
		expect += uint64(hashrate[i])
	}
	if tot := ethash.Hashrate(); tot != float64(expect) {
		t.Error("expect total hashrate should be same")
	}
}

func TestClosedRemoteSealer(t *testing.T) {
	ethash := NewTester(nil, false)
	time.Sleep(1 * time.Second) // ensure exit channel is listening
	ethash.Close()

	api := &API{ethash}
	if _, err := api.GetWork(); err != errEthashStopped {
		t.Error("expect to return an error to indicate ethash is stopped")
	}

	if res := api.SubmitHashRate(hexutil.Uint64(100), common.HexToHash("a")); res {
		t.Error("expect to return false when submit hashrate to a stopped ethash")
	}
}

// TestPrintDAGandCacheSizes isn't really a test; it's a just development tool to print scheduled
// DAG and cache sizes in human-readable form.
func TestPrintDAGandCacheSizes(t *testing.T) {
	logSizes := func(blockN uint64) {
		epoch := blockN / epochLength

		datasetS := datasetSize(epoch*epochLength + 1)
		cacheS := cacheSize(epoch*epochLength + 1)
		t.Logf("block=%d / epoch=%d -> data=%d (bin=%s dec=%s) cache=%d (bin=%s dec=%s)", blockN, blockN/epochLength,
			datasetS,
			byteCountBinary(datasetS),
			byteCountDecimal(datasetS),
			cacheS,
			byteCountBinary(cacheS),
			byteCountDecimal(cacheS),
		)
	}
	for _, n := range []uint64{0, 1, 2, 3} {
		logSizes(n * epochLength)
	}
	t.Log("---------------------------")
	for n := uint64(9_870_000); n <= 10_510_839; n += uint64(epochLength) {
		logSizes(n)
	}
	t.Log("---------------------------")
	for n := uint64(500_000); n < 16_000_000; n += uint64(500_000) {
		logSizes(n)
	}
}

// https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func byteCountDecimal(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func byteCountBinary(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// TestGenerateSizeStuntLogic tests the logic that stunts the dataset and cache size growth by epoch.
// It's a weak test because it just copy-pastes the logic and the logic is pretty simple logic anyways.
// But at least an I-wrote-the-test-and-want-to-be-sane kind of thing.
/*
// generate ensures that the dataset content is generated before use.
func (d *dataset) generate(dir string, limit int, test bool) {
	d.once.Do(func() {
		// Mark the dataset generated after we're done. This is needed for remote
		defer atomic.StoreUint32(&d.done, 1)

		// NOTE:ECIP1043
		var epoch = d.epoch
		if d.epochSizeStunt != 0 && d.epoch > d.epochSizeStunt {
			epoch = d.epochSizeStunt
		}
		csize := cacheSize(epoch*epochLength + 1)
		dsize := datasetSize(epoch*epochLength + 1)
*/
func TestGenerateSizeStuntLogic(t *testing.T) {
	testDataSet := []struct {
		epoch          uint64
		epochSizeStunt uint64
		wantSizeIndex  uint64
	}{
		{42, 40, 40},
		{38, 40, 38},
		{0, 0, 0},
		{0, 40, 0},
		{42, 0, 42},
	}
	for i, c := range testDataSet {

		// Reiterate the logic from the functions *dataset.generate and *cache.generate
		var epoch = c.epoch
		if c.epochSizeStunt != 0 && c.epoch > c.epochSizeStunt {
			epoch = c.epochSizeStunt
		}
		csize := cacheSize(epoch*epochLength + 1)
		dsize := datasetSize(epoch*epochLength + 1)

		// Make sure that the size fns returned with the expected values.
		if csize != cacheSizes[c.wantSizeIndex] {
			t.Error(i, "bad cache size", csize)
		}
		if dsize != datasetSizes[c.wantSizeIndex] {
			t.Error(i, "bad dataset size", dsize)
		}
	}
}
