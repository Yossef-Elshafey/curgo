package main

import (
	"log"
	"os"
)

func main() {
	path := "./examples/00.txt"
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file")
	}
	source := string(f)
	Interpret(source)
}
