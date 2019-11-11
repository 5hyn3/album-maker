package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(-1)
	}
}

var RootCmd = &cobra.Command{
	Use:   "culc",
	Short: "command line calculator",
	Run: func(cmd *cobra.Command, args []string) {
		var targetDir = "test-files/"

		files, err := ioutil.ReadDir(targetDir)
		if err != nil {
			panic(err)
		}

		var paths []string

		for _, file := range files {
			if !file.IsDir() {
				paths = append(paths, file.Name())
			}
		}

		for _, path := range paths {
			from := targetDir + path
			fileStat, err := os.Stat(from)
			if err != nil {
				log.Fatal(err)
			}
			var moveToDir = targetDir + fileStat.ModTime().Format("2006/01/02")
			if err := os.MkdirAll(moveToDir, 0777); err != nil {
				log.Fatal(err)
			}

			err = os.Rename(from, moveToDir+"/"+path)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
