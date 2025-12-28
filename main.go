package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.ReadFile("./examples/00.txt")
	if err != nil {
		log.Fatalf("Failed to read file")
	}
	source := string(f)
	Interpret(source)
}
