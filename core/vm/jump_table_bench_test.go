package vm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func BenchmarkJumpTable(b *testing.B) {
	c := params.DefaultMordorGenesisBlock().Config
	n := big.NewInt(3_000_000)
	for i := 0; i < b.N; i++ {
		instructionSetForConfig(c, n)
	}
}
