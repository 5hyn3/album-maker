package main

import (
	"github.com/5hyn3/album-maker/pkg/cmd"
	"github.com/5hyn3/album-maker/pkg/cmd/album-maker"
)

func main() {
	err := albummaker.NewCommand().Execute()
	cmd.CheckError(err)
}
