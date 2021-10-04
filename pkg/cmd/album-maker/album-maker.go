package albummaker

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/5hyn3/album-maker/pkg/cmd"
)

func MoveFilesToModTimeDirectory(targetDir string, mode SuffixMode) {
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

func moveFileToModTimeDirectory(targetDir string, path string, suffixMode SuffixMode, wg *sync.WaitGroup) {
	from := targetDir + "/" + path
	fileStat, err := os.Stat(from)
	cmd.CheckError(err)
	var moveToDir = targetDir + "/" + fileStat.ModTime().Format("2006/01/02")
	cmd.CheckError(os.MkdirAll(moveToDir, 0777))
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]
	suffix := ""
	switch suffixMode {
	case MD5:
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
	case DateTime:
		fileStat, err := os.Stat(from)
		if err != nil {
			panic(err)
		}
		suffix = "_" + fileStat.ModTime().Format("2006-01-02-15-04-05")
	default:
	}

	cmd.CheckError(os.Rename(from, moveToDir+"/"+base+suffix+ext))
	wg.Done()
}
