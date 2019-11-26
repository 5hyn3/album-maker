package albummaker

import (
	"os"
	"testing"
	"time"
)

var targetDirName = "test-files"

var testFileNames = []string{
	"test0",
	"test1",
	"test2",
}

func TestNewCommand(t *testing.T) {
	if err := setUp(); err != nil {
		panic(err)
	}

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

	if err := tearDown(); err != nil {
		panic(err)
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
