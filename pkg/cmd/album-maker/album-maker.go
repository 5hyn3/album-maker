package albummaker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/5hyn3/album-maker/internal/album-maker/entity"
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
	c.PersistentFlags().String(
		"suffixMode",
		"nothing",
		"Set suffix mode.",
		)
	return c
}

func makeAlbum(c *cobra.Command, args []string) {
	targetDir, err := c.PersistentFlags().GetString("targetDir")
	cmd.CheckError(err)

	suffixMode, err := c.PersistentFlags().GetString("suffixMode")
	cmd.CheckError(err)

	mode := entity.NewSuffixMode(suffixMode)

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
		go moveFileToModTimeDirectory(targetDir, path, mode, &wg)
	}
	wg.Wait()
}

func moveFileToModTimeDirectory(targetDir string, path string, suffixMode entity.SuffixMode, wg *sync.WaitGroup) {
	from := targetDir + "/" + path
	fileStat, err := os.Stat(from)
	cmd.CheckError(err)
	var moveToDir = targetDir + "/" + fileStat.ModTime().Format("2006/01/02")
	cmd.CheckError(os.MkdirAll(moveToDir, 0777))
	ext := filepath.Ext(path)
	base := path[:len(path) - len(ext)]
	suffix := ""
	switch suffixMode {
	case entity.MD5:
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
		suffix = "_" + hex.EncodeToString(h.Sum(nil))
	case entity.DateTime:
		fileStat, err := os.Stat(from)
		if err != nil {
			panic(err)
		}
		suffix = "_" + fileStat.ModTime().Format("2006-01-02-15-04-05")
	default:
	}

	cmd.CheckError(os.Rename(from, moveToDir + "/" + base + suffix + ext))
	wg.Done()
}
