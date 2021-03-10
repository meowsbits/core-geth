// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package ctypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// This file holds the Configurator interfaces.
// Interface methods follow distinct naming and signature patterns
// to enable abstracted logic.
//
// All methods are in pairs; Get'ers and Set'ers.
// Some Set methods are prefixed with Must, ie MustSet. These methods
// are allowed to return errors for debugging and logging, but
// any non-nil errors returned should stop program execution.
//
// All Forking methods (getting and setting Hard-fork requiring protocol changes)
// are suffixed with "Transition", and use *uint64 as in- and out-put variable types.
// A pointer is used because it's important that the fields can be nil, to signal
// being unset (as opposed to zero value uint64 == 0, ie from Genesis).

type Configurator interface {
	ChainConfigurator
	GenesisBlocker
}

type ChainConfigurator interface {
	ProtocolSpecifier
	Forker
	ConsensusEnginator // Consensus Engine
	// CHTer
}

// ProtocolSpecifier defines protocol interfaces that are agnostic of consensus engine.
type ProtocolSpecifier interface {
	GetAccountStartNonce() *uint64
	SetAccountStartNonce(n *uint64) error
	GetMaximumExtraDataSize() *uint64
	SetMaximumExtraDataSize(n *uint64) error
	GetMinGasLimit() *uint64
	SetMinGasLimit(n *uint64) error
	GetGasLimitBoundDivisor() *uint64
	SetGasLimitBoundDivisor(n *uint64) error
	GetNetworkID() *uint64
	SetNetworkID(n *uint64) error
	GetChainID() *big.Int
	SetChainID(i *big.Int) error
	GetMaxCodeSize() *uint64
	SetMaxCodeSize(n *uint64) error

	// Be careful with EIP2.
	// It is a messy EIP, specifying diverse changes, like difficulty, intrinsic gas costs for contract creation,
	// txpool management, and contract OoG handling.
	// It is both Ethash-specific and _not_.
	GetEIP2Transition() *uint64
	SetEIP2Transition(n *uint64) error

	GetEIP7Transition() *uint64
	SetEIP7Transition(n *uint64) error
	GetEIP150Transition() *uint64
	SetEIP150Transition(n *uint64) error
	GetEIP152Transition() *uint64
	SetEIP152Transition(n *uint64) error
	GetEIP160Transition() *uint64
	SetEIP160Transition(n *uint64) error
	GetEIP161abcTransition() *uint64
	SetEIP161abcTransition(n *uint64) error
	GetEIP161dTransition() *uint64
	SetEIP161dTransition(n *uint64) error
	GetEIP170Transition() *uint64
	SetEIP170Transition(n *uint64) error
	GetEIP155Transition() *uint64
	SetEIP155Transition(n *uint64) error
	GetEIP140Transition() *uint64
	SetEIP140Transition(n *uint64) error
	GetEIP198Transition() *uint64
	SetEIP198Transition(n *uint64) error
	GetEIP211Transition() *uint64
	SetEIP211Transition(n *uint64) error
	GetEIP212Transition() *uint64
	SetEIP212Transition(n *uint64) error
	GetEIP213Transition() *uint64
	SetEIP213Transition(n *uint64) error
	GetEIP214Transition() *uint64
	SetEIP214Transition(n *uint64) error
	GetEIP658Transition() *uint64
	SetEIP658Transition(n *uint64) error
	GetEIP145Transition() *uint64
	SetEIP145Transition(n *uint64) error
	GetEIP1014Transition() *uint64
	SetEIP1014Transition(n *uint64) error
	GetEIP1052Transition() *uint64
	SetEIP1052Transition(n *uint64) error
	GetEIP1283Transition() *uint64
	SetEIP1283Transition(n *uint64) error
	GetEIP1283DisableTransition() *uint64
	SetEIP1283DisableTransition(n *uint64) error
	GetEIP1108Transition() *uint64
	SetEIP1108Transition(n *uint64) error
	GetEIP2200Transition() *uint64
	SetEIP2200Transition(n *uint64) error
	GetEIP2200DisableTransition() *uint64
	SetEIP2200DisableTransition(n *uint64) error
	GetEIP1344Transition() *uint64
	SetEIP1344Transition(n *uint64) error
	GetEIP1884Transition() *uint64
	SetEIP1884Transition(n *uint64) error
	GetEIP2028Transition() *uint64
	SetEIP2028Transition(n *uint64) error
	GetECIP1080Transition() *uint64
	SetECIP1080Transition(n *uint64) error
	GetEIP1706Transition() *uint64
	SetEIP1706Transition(n *uint64) error
	GetEIP2537Transition() *uint64
	SetEIP2537Transition(n *uint64) error
	GetECBP1100Transition() *uint64
	SetECBP1100Transition(n *uint64) error
	GetEIP2315Transition() *uint64
	SetEIP2315Transition(n *uint64) error
	GetEIP2929Transition() *uint64
	SetEIP2929Transition(n *uint64) error
}

type Forker interface {
	// IsEnabled tells if interface has met or exceeded a fork block number.
	// eg. IsEnabled(c.GetEIP1108Transition, big.NewInt(42)))
	IsEnabled(fn func() *uint64, n *big.Int) bool

	// ForkCanonHash yields arbitrary number/hash pairs.
	// This is an abstraction derived from the original EIP150 implementation.
	GetForkCanonHash(n uint64) common.Hash
	SetForkCanonHash(n uint64, h common.Hash) error
	GetForkCanonHashes() map[uint64]common.Hash
}

type ConsensusEnginator interface {
	GetConsensusEngineType() ConsensusEngineT
	MustSetConsensusEngineType(t ConsensusEngineT) error
	EthashConfigurator
	CliqueConfigurator
}

type EthashConfigurator interface {
	GetEthashMinimumDifficulty() *big.Int
	SetEthashMinimumDifficulty(i *big.Int) error
	GetEthashDifficultyBoundDivisor() *big.Int
	SetEthashDifficultyBoundDivisor(i *big.Int) error
	GetEthashDurationLimit() *big.Int
	SetEthashDurationLimit(i *big.Int) error
	GetEthashHomesteadTransition() *uint64
	SetEthashHomesteadTransition(n *uint64) error

	// GetEthashEIP779Transition should return the block if the node wants the fork.
	// Otherwise, nil should be returned.
	GetEthashEIP779Transition() *uint64 // DAO

	// SetEthashEIP779Transition should turn DAO support on (nonnil) or off (nil).
	SetEthashEIP779Transition(n *uint64) error
	GetEthashEIP649Transition() *uint64
	SetEthashEIP649Transition(n *uint64) error
	GetEthashEIP1234Transition() *uint64
	SetEthashEIP1234Transition(n *uint64) error
	GetEthashEIP2384Transition() *uint64
	SetEthashEIP2384Transition(n *uint64) error
	GetEthashECIP1010PauseTransition() *uint64
	SetEthashECIP1010PauseTransition(n *uint64) error
	GetEthashECIP1010ContinueTransition() *uint64
	SetEthashECIP1010ContinueTransition(n *uint64) error
	GetEthashECIP1017Transition() *uint64
	SetEthashECIP1017Transition(n *uint64) error
	GetEthashECIP1017EraRounds() *uint64
	SetEthashECIP1017EraRounds(n *uint64) error
	GetEthashEIP100BTransition() *uint64
	SetEthashEIP100BTransition(n *uint64) error
	GetEthashECIP1041Transition() *uint64
	SetEthashECIP1041Transition(n *uint64) error
	GetEthashECIP1099Transition() *uint64
	SetEthashECIP1099Transition(n *uint64) error

	GetEthashDifficultyBombDelaySchedule() Uint64BigMapEncodesHex
	SetEthashDifficultyBombDelaySchedule(m Uint64BigMapEncodesHex) error
	GetEthashBlockRewardSchedule() Uint64BigMapEncodesHex
	SetEthashBlockRewardSchedule(m Uint64BigMapEncodesHex) error
}

type CliqueConfigurator interface {
	GetCliquePeriod() uint64
	SetCliquePeriod(n uint64) error
	GetCliqueEpoch() uint64
	SetCliqueEpoch(n uint64) error
}

type BlockSealer interface {
	GetSealingType() BlockSealingT
	SetSealingType(t BlockSealingT) error
	BlockSealerEthereum
}

type BlockSealerEthereum interface {
	GetGenesisSealerEthereumNonce() uint64
	SetGenesisSealerEthereumNonce(n uint64) error
	GetGenesisSealerEthereumMixHash() common.Hash
	SetGenesisSealerEthereumMixHash(h common.Hash) error
}

type GenesisBlocker interface {
	BlockSealer
	Accounter
	GetGenesisDifficulty() *big.Int
	SetGenesisDifficulty(i *big.Int) error
	GetGenesisAuthor() common.Address
	SetGenesisAuthor(a common.Address) error
	GetGenesisTimestamp() uint64
	SetGenesisTimestamp(u uint64) error
	GetGenesisParentHash() common.Hash
	SetGenesisParentHash(h common.Hash) error
	GetGenesisExtraData() []byte
	SetGenesisExtraData(b []byte) error
	GetGenesisGasLimit() uint64
	SetGenesisGasLimit(u uint64) error
}

type Accounter interface {
	ForEachAccount(fn func(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error) error
	UpdateAccount(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error
}

// normalizeUint64P normalizes the uint64 reference for use in ProtocolRules.
func normalizeUint64P(n *uint64) uint64 {
	if n == nil {
		return 0
	}
	return *n
}

// normalizeBigInt normalizes the big.Int reference for use in ProtocolRules.
func normalizeBigInt(n *big.Int) *big.Int {
	out := new(big.Int)
	return out.Set(n)
}

func GetProtocolRules(cc ChainConfigurator, blockNum *big.Int) ProtocolRules {
	return ProtocolRules{
		AccountStartNonce:    normalizeUint64P(cc.GetAccountStartNonce()),
		MaximumExtraDataSize: normalizeUint64P(cc.GetMaximumExtraDataSize()),
		MinGasLimit:          normalizeUint64P(cc.GetMinGasLimit()),
		GasLimitBoundDivisor: normalizeUint64P(cc.GetGasLimitBoundDivisor()),
		NetworkID:            normalizeUint64P(cc.GetNetworkID()),
		ChainID:              normalizeBigInt(cc.GetChainID()),
		MaxCodeSize:          normalizeUint64P(cc.GetMaxCodeSize()),

		EIP2Transition_Enabled:           cc.IsEnabled(cc.GetEIP2Transition, blockNum),
		EIP7Transition_Enabled:           cc.IsEnabled(cc.GetEIP7Transition, blockNum),
		EIP150Transition_Enabled:         cc.IsEnabled(cc.GetEIP150Transition, blockNum),
		EIP152Transition_Enabled:         cc.IsEnabled(cc.GetEIP152Transition, blockNum),
		EIP160Transition_Enabled:         cc.IsEnabled(cc.GetEIP160Transition, blockNum),
		EIP161abcTransition_Enabled:      cc.IsEnabled(cc.GetEIP161abcTransition, blockNum),
		EIP161dTransition_Enabled:        cc.IsEnabled(cc.GetEIP161dTransition, blockNum),
		EIP170Transition_Enabled:         cc.IsEnabled(cc.GetEIP170Transition, blockNum),
		EIP155Transition_Enabled:         cc.IsEnabled(cc.GetEIP155Transition, blockNum),
		EIP140Transition_Enabled:         cc.IsEnabled(cc.GetEIP140Transition, blockNum),
		EIP198Transition_Enabled:         cc.IsEnabled(cc.GetEIP198Transition, blockNum),
		EIP211Transition_Enabled:         cc.IsEnabled(cc.GetEIP211Transition, blockNum),
		EIP212Transition_Enabled:         cc.IsEnabled(cc.GetEIP212Transition, blockNum),
		EIP213Transition_Enabled:         cc.IsEnabled(cc.GetEIP213Transition, blockNum),
		EIP214Transition_Enabled:         cc.IsEnabled(cc.GetEIP214Transition, blockNum),
		EIP658Transition_Enabled:         cc.IsEnabled(cc.GetEIP658Transition, blockNum),
		EIP145Transition_Enabled:         cc.IsEnabled(cc.GetEIP145Transition, blockNum),
		EIP1014Transition_Enabled:        cc.IsEnabled(cc.GetEIP1014Transition, blockNum),
		EIP1052Transition_Enabled:        cc.IsEnabled(cc.GetEIP1052Transition, blockNum),
		EIP1283Transition_Enabled:        cc.IsEnabled(cc.GetEIP1283Transition, blockNum),
		EIP1283DisableTransition_Enabled: cc.IsEnabled(cc.GetEIP1283DisableTransition, blockNum),
		EIP1108Transition_Enabled:        cc.IsEnabled(cc.GetEIP1108Transition, blockNum),
		EIP2200Transition_Enabled:        cc.IsEnabled(cc.GetEIP2200Transition, blockNum),
		EIP2200DisableTransition_Enabled: cc.IsEnabled(cc.GetEIP2200DisableTransition, blockNum),
		EIP1344Transition_Enabled:        cc.IsEnabled(cc.GetEIP1344Transition, blockNum),
		EIP1884Transition_Enabled:        cc.IsEnabled(cc.GetEIP1884Transition, blockNum),
		EIP2028Transition_Enabled:        cc.IsEnabled(cc.GetEIP2028Transition, blockNum),
		ECIP1080Transition_Enabled:       cc.IsEnabled(cc.GetECIP1080Transition, blockNum),
		EIP1706Transition_Enabled:        cc.IsEnabled(cc.GetEIP1706Transition, blockNum),
		EIP2537Transition_Enabled:        cc.IsEnabled(cc.GetEIP2537Transition, blockNum),
		ECBP1100Transition_Enabled:       cc.IsEnabled(cc.GetECBP1100Transition, blockNum),
		EIP2315Transition_Enabled:        cc.IsEnabled(cc.GetEIP2315Transition, blockNum),
		EIP2929Transition_Enabled:        cc.IsEnabled(cc.GetEIP2929Transition, blockNum),
	}
}

type ProtocolRules struct {
	AccountStartNonce    uint64
	MaximumExtraDataSize uint64
	MinGasLimit          uint64
	GasLimitBoundDivisor uint64
	NetworkID            uint64
	ChainID              *big.Int
	MaxCodeSize          uint64

	// Be careful with EIP2.
	// It is a messy EIP, specifying diverse changes, like difficulty, intrinsic gas costs for contract creation,
	// txpool management, and contract OoG handling.
	// It is both Ethash-specific and _not_.
	EIP2Transition_Enabled           bool
	EIP7Transition_Enabled           bool
	EIP150Transition_Enabled         bool
	EIP152Transition_Enabled         bool
	EIP160Transition_Enabled         bool
	EIP161abcTransition_Enabled      bool
	EIP161dTransition_Enabled        bool
	EIP170Transition_Enabled         bool
	EIP155Transition_Enabled         bool
	EIP140Transition_Enabled         bool
	EIP198Transition_Enabled         bool
	EIP211Transition_Enabled         bool
	EIP212Transition_Enabled         bool
	EIP213Transition_Enabled         bool
	EIP214Transition_Enabled         bool
	EIP658Transition_Enabled         bool
	EIP145Transition_Enabled         bool
	EIP1014Transition_Enabled        bool
	EIP1052Transition_Enabled        bool
	EIP1283Transition_Enabled        bool
	EIP1283DisableTransition_Enabled bool
	EIP1108Transition_Enabled        bool
	EIP2200Transition_Enabled        bool
	EIP2200DisableTransition_Enabled bool
	EIP1344Transition_Enabled        bool
	EIP1884Transition_Enabled        bool
	EIP2028Transition_Enabled        bool
	ECIP1080Transition_Enabled       bool
	EIP1706Transition_Enabled        bool
	EIP2537Transition_Enabled        bool
	ECBP1100Transition_Enabled       bool
	EIP2315Transition_Enabled        bool
	EIP2929Transition_Enabled        bool
}
