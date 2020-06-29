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
	var g *genesisT.Genesis
	err = json.Unmarshal(b, g)
	if err != nil {
		t.Fatal(err)
	}
	
}
