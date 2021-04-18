// Copyright 2020 The go-ethereum Authors
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

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// AccessListSegmentIDTx is the data of the iip999 transaction type.
// This transaction type extends EIP2929's optional access list pattern,
// including a segment ID field.
type AccessListSegmentIDTx struct {
	ChainID    *big.Int        // destination chain ID
	SegmentID  *big.Int        // destination segment ID
	Nonce      uint64          // nonce of sender account
	GasPrice   *big.Int        // wei per gas
	Gas        uint64          // gas limit
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int        // wei amount
	Data       []byte          // contract invocation input data
	AccessList AccessList      // EIP-2930 access list
	V, R, S    *big.Int        // signature values
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *AccessListSegmentIDTx) copy() TxData {
	cpy := &AccessListSegmentIDTx{
		Nonce: tx.Nonce,
		To:    tx.To, // TODO: copy pointed-to address
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		SegmentID:  new(big.Int),
		GasPrice:   new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	copy(cpy.AccessList, tx.AccessList)
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.SegmentID != nil {
		cpy.SegmentID.Set(tx.SegmentID)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.

func (tx *AccessListSegmentIDTx) txType() byte           { return AccessListSegmentIDTxType }
func (tx *AccessListSegmentIDTx) chainID() *big.Int      { return tx.ChainID }
func (tx *AccessListSegmentIDTx) segmentID() *big.Int    { return tx.SegmentID }
func (tx *AccessListSegmentIDTx) protected() bool        { return true }
func (tx *AccessListSegmentIDTx) accessList() AccessList { return tx.AccessList }
func (tx *AccessListSegmentIDTx) data() []byte           { return tx.Data }
func (tx *AccessListSegmentIDTx) gas() uint64            { return tx.Gas }
func (tx *AccessListSegmentIDTx) gasPrice() *big.Int     { return tx.GasPrice }
func (tx *AccessListSegmentIDTx) value() *big.Int        { return tx.Value }
func (tx *AccessListSegmentIDTx) nonce() uint64          { return tx.Nonce }
func (tx *AccessListSegmentIDTx) to() *common.Address    { return tx.To }

func (tx *AccessListSegmentIDTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *AccessListSegmentIDTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}
