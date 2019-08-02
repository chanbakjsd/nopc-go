package ast_test

import (
	"path/filepath"
	"testing"

	"github.com/noplang/nopc-go/ast"
)

func TestHelloWorld(t *testing.T) {
	testFile("helloworld", false, t)
}

func TestDiscordDiscussion(t *testing.T) {
	//Test that AST generation works for the sample (theoratical) code discussed in Discord.
	testFile("discord_discussion", false, t)
}

func TestComplexExample(t *testing.T) {
	testFile("complex", false, t)
}

func TestFailLexer(t *testing.T) {
	testFile("lexer_fail", true, t)
}

func TestFailParser(t *testing.T) {
	testFile("parser_fail", true, t)
}

func testFile(testFileName string, shouldFail bool, test *testing.T) {
	testFileLocation := filepath.Join("testdata", testFileName+".nop")

	tree := ast.ParseFile(testFileLocation)
	if tree == nil {
		if !shouldFail {
			test.Fatal("An error has occurred while parsing a file")
		}
		return
	}
	if tree != nil && shouldFail {
		test.Fatal("File that should fail did not fail!")
	}
	generatedAST := ast.GenerateAST([]ast.ParseTreeWithLocation{*tree})
	prettyPrinted := ast.PrettyPrintAST(generatedAST[0])

	if err := runTestOnFile(prettyPrinted, testFileName+".ast.golden"); err != nil {
		test.Error(err)
	}
}
