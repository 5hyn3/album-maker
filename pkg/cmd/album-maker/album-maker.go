package albummaker

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MoveFilesToModTimeDirectory(targetDir string, mode SuffixMode) error {
	files, err := ioutil.ReadDir(targetDir)

	if err != nil {
		return err
	}

	var paths []string

	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, file.Name())
		}
	}

	eg := errgroup.Group{}
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			return moveFileToModTimeDirectory(targetDir, path, mode)
		})
	}
	err = eg.Wait()
	if  err != nil {
		return err
	}
	return nil
}

func moveFileToModTimeDirectory(targetDir string, path string, suffixMode SuffixMode) error {
	from := targetDir + "/" + path
	fileStat, err := os.Stat(from)
	if err != nil {
		return err
	}

	var moveToDir = targetDir + "/" + fileStat.ModTime().Format("2006/01/02")
	err = os.MkdirAll(moveToDir, 0777)
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]
	suffix := ""
	switch suffixMode {
	case MD5:
		f, err := os.Open(from)
		if err != nil {
			return err
		}

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			panic(err)
		}

		err = f.Close()
		if err != nil {
			return err
		}
		suffix = "_" + hex.EncodeToString(h.Sum(nil))
	case DateTime:
		fileStat, err := os.Stat(from)
		if err != nil {
			return err
		}
		suffix = "_" + fileStat.ModTime().Format("2006-01-02-15-04-05")
	default:
	}

	err = os.Rename(from, moveToDir+"/"+base+suffix+ext)
	if err != nil {
		return err
	}
	return nil
}