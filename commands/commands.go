package commands

import (
	"curgo/lexer"
	"curgo/storage"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CommandHandler struct {
	Fs       *flag.FlagSet
	cf       string
	init     bool
	filename string
	block    string
}

func NewCommandHandler() *CommandHandler {
	ch := &CommandHandler{Fs: flag.NewFlagSet(os.Args[0], flag.ContinueOnError)}
	ch.Fs.StringVar(&ch.cf, "createf", "", "create a new file (without extension)")
	ch.Fs.BoolVar(&ch.init, "init", false, "create curgo")
	ch.Fs.StringVar(&ch.filename, "f", "", "points to a file")
	ch.Fs.StringVar(&ch.block, "c", "", "points to a block")
	return ch
}

func (ch *CommandHandler) CreateFileFn() {
	if ch.cf != "" {
		fmt.Println(storage.CreateFile(ch.cf))
	}
}

func (ch *CommandHandler) InitFn() {
	if ch.init {
		storage.CreateDir(storage.ROOT)
		storage.CreateFile("root")
	}
}

func (ch *CommandHandler) ExecuteFullFile() {
	if ch.block != "" && ch.filename != "" {
		content := storage.ReadFile(ch.filename)
		l := lexer.NewLexer(content)
		block, err := strconv.Atoi(ch.block)
		if err != nil {
			fmt.Printf("Cannot process block:%s\n", ch.block)
			os.Exit(1)
		}
		l.Tokenize(block)
	}
}

func (ch *CommandHandler) ExecuteBlock() {
	if ch.block == "" && ch.filename != "" {
		content := storage.ReadFile(ch.filename)
		l := lexer.NewLexer(content)
		l.Tokenize(-1)
	}
}
