package main

import (
	"curgo/parser"
	"os"
)

func main() {
	parser.NewParser()
	initCommands()
}

func initCommands() {
	ch := NewCommandHandler()
	ch.fs.Parse(os.Args[1:])
	ch.createFileFn()
	ch.initFn()
}
