/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/spf13/cobra"
)

// writeHeadblockhashCmd represents the writeHeadblockhash command
var writeHeadblockhashCmd = &cobra.Command{
	Use:   "write-headblockhash",
	Short: "Write head block hash to the chain database",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("need hash")
		}

		log.Println("Opening database...")
		db, err := rawdb.NewLevelDBDatabase(chainDBPath, 256, 16, "")
		if err != nil {
			log.Fatal(err)
		}

		rawdb.WriteHeadBlockHash(db, common.HexToHash(args[1]))
		i, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		rawdb.WriteCanonicalHash(db, common.HexToHash(args[1]), i)
	},
}

func init() {
	rootCmd.AddCommand(writeHeadblockhashCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeHeadblockhashCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeHeadblockhashCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
