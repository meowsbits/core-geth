package integration

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// The following benchmarks compare ways of getting
// configuration values from a configuration struct.

// Example result:
/*
goos: linux
goarch: amd64
pkg: github.com/ethereum/go-ethereum/integration
BenchmarkConfigurationRules
BenchmarkConfigurationRules-12        	1000000000	         0.423 ns/op
BenchmarkConfigurationInterface
BenchmarkConfigurationInterface-12    	16766822	        64.6 ns/op
PASS
*/

func BenchmarkConfigurationRules(b *testing.B) {
	g := params.DefaultMordorGenesisBlock()
	rules := ctypes.GetProtocolRules(g, big.NewInt(3000000))
	for i := 0; i < b.N; i++ {
		if rules.EIP161abcTransition_Enabled {
		}
	}
}

func BenchmarkConfigurationInterface(b *testing.B) {
	g := params.DefaultMordorGenesisBlock()
	n := big.NewInt(3000000)
	for i := 0; i < b.N; i++ {
		if g.IsEnabled(g.GetEIP161abcTransition, n) {
		}
	}
}
