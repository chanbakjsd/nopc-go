package main

import (
	"./flags"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func main() {
	files := flag.Args()
	if len(files) == 0 {
		panic("REPL has not been implemented yet.")
	}

	parsedTree := make([]antlr.Tree, len(files))
	encounteredError := false
	for i := 0; i < len(files); i++ {
		flags.DebugPrint(files[i] + ": processing")
		parsedTree[i] = parseFile(files[i])
		if parsedTree[i] == nil {
			encounteredError = true
			flags.DebugPrint(files[i] + ": process failed, proceeding to other files")
		} else {
			flags.DebugPrint(files[i] + ": processed")
		}
	}
	if encounteredError {
		fmt.Println("Error(s) have occurred while processing source files.")
		os.Exit(1)
	}
	if flags.OutputLocation == "<auto>" {
		flags.OutputLocation = strings.TrimSuffix(files[0], ".nop")
	}
	flags.DebugPrint("Output file location: " + flags.OutputLocation)
}
