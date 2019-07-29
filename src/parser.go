package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/noplang/nopc-go/ast"
	"github.com/noplang/nopc-go/flags"
	"github.com/noplang/nopc-go/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type parseTreeWithLocation struct {
	Tree     antlr.ParseTree
	Location string
}

func parseFile(location string) *parseTreeWithLocation {
	input, err := antlr.NewFileStream(location)
	if err != nil {
		fmt.Println(location + ": Error parsing file through antlr: " + err.Error())
		return nil
	}

	lexerError := errorListener{}
	parserError := errorListener{}

	lexer := parser.NewnopLexer(input)
	lexer.AddErrorListener(&lexerError)

	stream := antlr.NewCommonTokenStream(lexer, 0)

	parser := parser.NewnopParser(stream)
	parser.BuildParseTrees = true
	parser.AddErrorListener(&parserError)
	tree := parser.Nop_file()

	if lexerError.EncounteredError {
		fmt.Println(location + ": Failed to tokenize file")
		return nil
	}
	if parserError.EncounteredError {
		fmt.Println(location + ": Failed to parse file")
		return nil
	}
	return &parseTreeWithLocation{
		Tree:     tree,
		Location: location,
	}
}

func parseFiles(location []string) ([]parseTreeWithLocation, bool) {
	var encounteredError bool
	result := make([]parseTreeWithLocation, 0, len(location))
	for i := 0; i < len(location); i++ {
		err := filepath.Walk(location[i], func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil //No need to parse a directory
			}
			if !strings.HasSuffix(path, ".nop") {
				return nil //Nop file only
			}
			flags.DebugPrint(path + ": processing")
			nextTree := parseFile(path)
			if nextTree == nil {
				flags.DebugPrint(path + ": skipping file due to error")
				encounteredError = true
				return nil
			}
			result = append(result, *nextTree)
			flags.DebugPrint(path + ": processed")

			return nil
		})
		if err != nil {
			fmt.Println(location[i] + ": walk error: " + err.Error())
			encounteredError = true
		}
	}
	return result, encounteredError
}

func generateAST(parseTrees []parseTreeWithLocation) []ast.NopFile {
	result := make([]ast.NopFile, 0, len(parseTrees))
	for i, v := range parseTrees {
		flags.DebugPrint(v.Location + ": walking tree")
		astWalker := ast.NewNopListener()
		nopTree := astWalker.Walk(parseTrees[i].Tree)
		nopTree.Source = parseTrees[i].Location
		result = append(result, nopTree)
		flags.DebugPrint(v.Location + ": walked tree")
	}
	return result
}
