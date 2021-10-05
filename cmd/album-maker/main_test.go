package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var targetDirName = "test-files"

var testFileNames = []string{
	"test0.txt",
	"test1.txt",
	"test2.txt",
}

func TestNotPassedSuffixModePattern(t *testing.T) {
	if err := setUp(); err != nil {
		panic(err)
	}
	defer tearDown()

	var args = []string{"--targetDir=" + targetDirName}
	var cmd = NewCommand()
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	var todayDir = time.Now().Format("2006/01/02")
	for _, n := range testFileNames {
		var target = targetDirName + "/" + todayDir + "/" + n
		if !fileExists(target) {
			t.Fatalf("FAIL: %s\n", target)
		}
	}
}

func TestSuffixModeEqualNothingPattern(t *testing.T) {
	if err := setUp(); err != nil {
		panic(err)
	}
	defer tearDown()

	var args = []string{"--targetDir=" + targetDirName, "--suffixMode=nothing"}
	var cmd = NewCommand()
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	var todayDir = time.Now().Format("2006/01/02")
	for _, n := range testFileNames {
		var target = targetDirName + "/" + todayDir + "/" + n
		if !fileExists(target) {
			t.Fatalf("FAIL: %s\n", target)
		}
	}
}

func TestSuffixModeEqualMd5Pattern(t *testing.T) {
	if err := setUp(); err != nil {
		panic(err)
	}
	defer tearDown()

	var nameToMd5 = map[string]string{}
	for _, n := range testFileNames {
		target := targetDirName + "/" + n
		f, err := os.Open(target)
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

		hash := hex.EncodeToString(h.Sum(nil))
		nameToMd5[n] = hash
	}

	var args = []string{"--targetDir=" + targetDirName, "--suffixMode=md5"}
	var cmd = NewCommand()
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	var todayDir = time.Now().Format("2006/01/02")
	for _, n := range testFileNames {
		ext := filepath.Ext(n)
		base := n[:len(n)-len(ext)]
		var target = targetDirName + "/" + todayDir + "/" + base + "_" + nameToMd5[n] + ext
		if !fileExists(target) {
			t.Fatalf("FAIL: %s\n", target)
		}
	}
}

func TestSuffixModeEqualDateTimePattern(t *testing.T) {
	if err := setUp(); err != nil {
		panic(err)
	}
	defer tearDown()

	var nameToDateTime = map[string]string{}
	for _, n := range testFileNames {
		target := targetDirName + "/" + n
		fileStat, err := os.Stat(target)
		if err != nil {
			panic(err)
		}
		var dateTime = fileStat.ModTime().Format("2006-01-02-15-04-05")

		nameToDateTime[n] = dateTime
	}

	var args = []string{"--targetDir=" + targetDirName, "--suffixMode=dateTime"}
	var cmd = NewCommand()
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	var todayDir = time.Now().Format("2006/01/02")
	for _, n := range testFileNames {
		ext := filepath.Ext(n)
		base := n[:len(n)-len(ext)]
		dateTime := nameToDateTime[n]
		var target = targetDirName + "/" + todayDir + "/" + base + "_" + dateTime + ext
		if !fileExists(target) {
			t.Fatalf("FAIL: %s\n", target)
		}
	}
}

func setUp() error {
	err := os.MkdirAll(targetDirName, 0777)
	if err != nil {
		return err
	}
	for _, n := range testFileNames {
		file, err := os.Create(targetDirName + `/` + n)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func tearDown() error {
	return os.RemoveAll(`test-files`)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
