package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestMigrateFaucetDirectory(t *testing.T) {
	hardToCollideName := fmt.Sprintf("faucet-migration-test-datadir-%d", time.Now().UnixNano())
	tempDir := filepath.Join(os.TempDir(), hardToCollideName)
	defer func() {
		os.RemoveAll(tempDir)
	}()

	faucetDataDir := filepath.Join(tempDir, "mychain")
	oldFaucetNodeDataDir := filepath.Join(faucetDataDir, multiFaucetNodeName)

	if err := os.MkdirAll(oldFaucetNodeDataDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	filepath.Walk(faucetDataDir, func(path string, info os.FileInfo, err error) error {
		t.Logf("%s", path)
		return nil
	})

	if err := migrateFaucetDirectory(faucetDataDir); err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(faucetDataDir, coreFaucetNodeName)
	d, err := os.Stat(expected)
	if err != nil {
		t.Fatal(err)
	}
	if !d.IsDir() {
		t.Fatal("non-directory")
	}
	filepath.Walk(faucetDataDir, func(path string, info os.FileInfo, err error) error {
		t.Logf("%s", path)
		return nil
	})
}

// TestFacebook makes a live request to Facebook to retrieve and parse
// an example post.
// This test consistently fails on Travis CI (possibly because a blacklist either at Travis or Facebook).
func TestFacebook(t *testing.T) {
	for _, tt := range []struct {
		url  string
		want common.Address
	}{
		{
			"https://www.facebook.com/fooz.gazonk/posts/2837228539847129",
			common.HexToAddress("0xDeadDeaDDeaDbEefbEeFbEEfBeeFBeefBeeFbEEF"),
		},
	} {
		_, _, gotAddress, err := authFacebook(tt.url)
		if err != nil {
			t.Fatal(err)
		}
		if gotAddress != tt.want {
			t.Fatalf("address wrong, have %v want %v", gotAddress, tt.want)
		}
	}
}
