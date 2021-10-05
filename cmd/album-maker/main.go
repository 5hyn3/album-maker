package main

import (
	"fmt"

	"github.com/5hyn3/album-maker/pkg/cmd"
	albummaker "github.com/5hyn3/album-maker/pkg/cmd/album-maker"
	"github.com/spf13/cobra"
)

func main() {
	err := NewCommand().Execute()
	cmd.CheckError(err)
}

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

	mode := albummaker.NewSuffixMode(suffixMode)

	if targetDir == "" {
		fmt.Print("TargetDir must be set.")
		return
	}

	err = albummaker.MoveFilesToModTimeDirectory(targetDir, mode)
	cmd.CheckError(err)
}
