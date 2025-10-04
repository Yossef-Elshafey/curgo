package main

import (
	"curgo/commands"
	"os"
)

func main() {
	initCommands()
}

func initCommands() {
	ch := commands.NewCommandHandler()
	ch.Fs.Parse(os.Args[1:])
	ch.CreateFileFn()
	ch.InitFn()
	ch.ExecuteFullFile()
	ch.ExecuteBlock()
}
