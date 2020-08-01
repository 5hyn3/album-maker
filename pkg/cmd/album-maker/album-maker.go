package albummaker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/5hyn3/album-maker/pkg/cmd"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func NewCommand() *cobra.Command {
	var c = &cobra.Command{
		Use: "album maker",
		// ファイル（主に画像）を最終変更日を元にディレクトリ分け整理します。
		Short: "Organize files (mainly images) by dividing them into directories based on the last modification date.",
		Run:   makeAlbum,
	}
	c.PersistentFlags().String("targetDir", "", "Set target directory.")
	c.PersistentFlags().Bool(
		"addMd5Suffix",
		false,
		"Set the suffix of the original file name to include the md5 calculated from the file.",
		)
	return c
}

func makeAlbum(c *cobra.Command, args []string) {
	targetDir, err := c.PersistentFlags().GetString("targetDir")
	cmd.CheckError(err)

	addMd5Suffix, err := c.PersistentFlags().GetBool("addMd5Suffix")
	cmd.CheckError(err)

	if targetDir == "" {
		fmt.Print("TargetDir must be set.")
		return
	}

	files, err := ioutil.ReadDir(targetDir)
	cmd.CheckError(err)

	var paths []string

	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, file.Name())
		}
	}

	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go moveFileToModTimeDirectory(targetDir, path, addMd5Suffix, &wg)
	}
	wg.Wait()
}

func moveFileToModTimeDirectory(targetDir string, path string, addMd5Suffix bool, wg *sync.WaitGroup) {
	from := targetDir + "/" + path
	fileStat, err := os.Stat(from)
	cmd.CheckError(err)
	var moveToDir = targetDir + "/" + fileStat.ModTime().Format("2006/01/02")
	cmd.CheckError(os.MkdirAll(moveToDir, 0777))
	var newName string
	if addMd5Suffix {
		f, err := os.Open(from)
		if err != nil {
			panic(err)
		}

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			panic(err)
		}

		err = f.Close()
		if err != nil {
			panic(err)
		}

		ext := filepath.Ext(path)
		base := path[:len(path) - len(ext)]
		newName = moveToDir + "/" + base + "_" + hex.EncodeToString(h.Sum(nil)) + ext
	} else {
		newName = moveToDir + "/" + path
	}
	cmd.CheckError(os.Rename(from, newName))
	wg.Done()
}
