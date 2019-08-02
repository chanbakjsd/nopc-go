package ast_test

import (
	"testing"

	"github.com/noplang/nopc-go/ast"
)

func TestCleanNopFile(t *testing.T) {
	output := ast.PrettyPrintAST(ast.NopFile{})
	err := runTestOnFile(output, "pretty_print_clean.json.golden")
	if err != nil {
		t.Error(err)
	}
}
