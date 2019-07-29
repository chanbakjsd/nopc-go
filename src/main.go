package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/noplang/nopc-go/ast"
	"github.com/noplang/nopc-go/flags"
)

func main() {
	files := flag.Args()
	if len(files) == 0 {
		panic("REPL has not been implemented yet.")
	}

	parsedTree, encounteredError := parseFiles(files)
	if encounteredError {
		fmt.Println("Error(s) have occurred while processing source files.")
		os.Exit(1)
	}
	if flags.OutputLocation == "<auto>" {
		flags.OutputLocation = strings.TrimSuffix(files[0], ".nop")
	}
	flags.DebugPrint("Output file location: " + flags.OutputLocation)
	astTree := generateAST(parsedTree)
	if flags.EmitAST {
		for _, v := range astTree {
			ast.PrettyPrintASTToFile(v)
		}
		fmt.Println("AST files can be found with their source (.ast)")
		return
	}
	fmt.Println("Nop compiled successfully.")
}
