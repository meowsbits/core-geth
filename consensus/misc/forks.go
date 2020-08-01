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

package misc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// VerifyForkHashes verifies that blocks conforming to network hard-forks do have
// the correct hashes, to avoid clients going off on different chains. This is an
// optional feature.
func VerifyForkHashes(config ctypes.ChainConfigurator, header *types.Header, uncle bool) error {
	// We don't care about uncles
	if uncle {
		return nil
	}

	// Development debug:
	// Add the canonical chain hash as a disallowed block to hopefully
	// force the chain to sync with the reorged (non-canonical chain)
	// block: 10904147 (<- first reorg number, then continuing ~3500 blocks).
	canonicalNumber := uint64(10904147)
	canonicalHash := common.HexToHash("0x7a875db2878a4e7ca63357f5fbcce4fce47d17e15c116ac04ab20aaaab9dbca6")
	if header.Number.Uint64() == canonicalNumber && header.Hash() == canonicalHash {
		log.Warn("Selecting reorged chain segment intentionally")
		j, _ := header.MarshalJSON()
		log.Warn("Selecting reorged header: %s", string(j))
		return fmt.Errorf("DEBUG: ETC chainsplit reorged fork selection: number=%d canonical=%s using=%s", canonicalNumber, canonicalHash.Hex(), header.Hash().Hex())
	}
	if wantHash := config.GetForkCanonHash(header.Number.Uint64()); wantHash == (common.Hash{}) || wantHash == header.Hash() {
		return nil
	} else {
		return fmt.Errorf("verify canonical block hash failed, block number %d: have 0x%x, want 0x%x", header.Number.Uint64(), header.Hash(), wantHash)
	}
}
