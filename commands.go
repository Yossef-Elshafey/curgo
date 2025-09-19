package main

import (
	"curgo/parser/storage"
	"flag"
	"fmt"
	"os"
)

type CommandHandler struct {
	fs   *flag.FlagSet
	cf   string
	init bool
}

func NewCommandHandler() *CommandHandler {
	ch := &CommandHandler{fs: flag.NewFlagSet(os.Args[0], flag.ContinueOnError)}
	ch.fs.StringVar(&ch.cf, "createf", "", "create a new file (without extension)")
	ch.fs.BoolVar(&ch.init, "init", false, "create curgo")
	return ch
}

func (ch *CommandHandler) createFileFn() {
	if ch.cf != "" {
		fmt.Println(storage.CreateFile(ch.cf))
	}
}

func (ch *CommandHandler) initFn() {
	if ch.init {
		storage.CreateDir(storage.ROOT)
		storage.CreateFile("root")
	}
}
