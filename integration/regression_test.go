package integration

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

func TestConfig1(t *testing.T) {
	fp := "../foundation.conf.json"
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}
	var g = new(genesisT.Genesis)
	err = json.Unmarshal(b, g)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGas(t *testing.T) {
	msgGas := uint64(9214364837600034817)
	lim := uint64(10000000)
	if msgGas > lim {
		msgGas = lim
	}
	t.Log(msgGas)
}