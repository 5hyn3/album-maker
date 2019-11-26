package albummaker

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func NewCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: "album maker",
		// ファイル（主に画像）を最終変更日を元にディレクトリ分け整理します。
		Short: "Organize files (mainly images) by dividing them into directories based on the last modification date.",
		Run: func(cmd *cobra.Command, args []string) {

			targetDir, err := cmd.PersistentFlags().GetString("targetDir")
			if err != nil {
				log.Fatal(err)
			}

			if targetDir == "" {
				fmt.Print("TargetDir must be set.")
				return
			}

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

			var wg sync.WaitGroup
			for _, path := range paths {
				wg.Add(1)
				go moveFileToModTimeDirectory(targetDir, path, &wg)
			}
			wg.Wait()
		},
	}
	rootCmd.PersistentFlags().String("targetDir", "", "Set target directory.")
	return rootCmd
}

func moveFileToModTimeDirectory(targetDir string, path string, wg *sync.WaitGroup) {
	from := targetDir + "/" + path
	fileStat, err := os.Stat(from)
	if err != nil {
		log.Fatal(err)
	}
	var moveToDir = targetDir + "/" + fileStat.ModTime().Format("2006/01/02")
	if err := os.MkdirAll(moveToDir, 0777); err != nil {
		log.Fatal(err)
	}
	err = os.Rename(from, moveToDir + "/" + path)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
}