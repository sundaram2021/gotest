// main.go
package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/sundaram2021/gotest/generator"
)

func main() {
	// Parse command-line arguments
	var inputFile string
	flag.StringVar(&inputFile, "file", "", "Path to the Go source file")
	flag.Parse()

	if inputFile == "" {
		fmt.Println("Please provide a Go source file using the -file flag.")
		os.Exit(1)
	}

	// Lexical analysis and parsing
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, inputFile, nil, 0)
	if err != nil {
		fmt.Printf("Failed to parse file: %v\n", err)
		os.Exit(1)
	}

	// Generate tests
	generator.GenerateTests(node, inputFile)
}
