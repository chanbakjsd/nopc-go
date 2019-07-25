package main

import (
	"./parser"
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func parseFile(location string) antlr.Tree {
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
	return tree
}
