package ast

import (
	"fmt"

	"github.com/noplang/nopc-go/flags"
	"github.com/noplang/nopc-go/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

//ParseTreeWithLocation : A struct that contains Antlr's parse tree and the location of the original file for debugging.
type ParseTreeWithLocation struct {
	Tree     antlr.ParseTree
	Location string
}

//ParseFile : A function that takes in the location of a file and outputs the Antlr parse tree with its location.
func ParseFile(location string) *ParseTreeWithLocation {
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
	return &ParseTreeWithLocation{
		Tree:     tree,
		Location: location,
	}
}

//GenerateAST : A function that takes in a list of Antlr parse trees and generate a list of AST.
func GenerateAST(parseTrees []ParseTreeWithLocation) []NopFile {
	result := make([]NopFile, 0, len(parseTrees))
	for i, v := range parseTrees {
		flags.DebugPrint(v.Location + ": walking tree")
		astWalker := NewNopListener()
		nopTree := astWalker.Walk(parseTrees[i].Tree)
		nopTree.Source = parseTrees[i].Location
		result = append(result, nopTree)
		flags.DebugPrint(v.Location + ": walked tree")
	}
	return result
}
