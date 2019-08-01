package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/noplang/nopc-go/ast"
	"github.com/noplang/nopc-go/flags"
)

func parseFiles(location []string) ([]ast.ParseTreeWithLocation, bool) {
	var encounteredError bool
	result := make([]ast.ParseTreeWithLocation, 0, len(location))
	for i := 0; i < len(location); i++ {
		err := filepath.Walk(location[i], func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil //No need to parse a directory
			}
			if !strings.HasSuffix(path, ".nop") {
				return nil //Nop file only
			}
			flags.DebugPrint(path + ": processing")
			nextTree := ast.ParseFile(path)
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
