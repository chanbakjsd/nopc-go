package flags

import "flag"

var (
	OutputLocation string
	Debug          bool
	EmitAST        bool
)

func init() {
	flag.StringVar(&OutputLocation, "output", "<auto>", "Where to output the compiled executable")
	flag.StringVar(&OutputLocation, "o", "<auto>", "Where to output the compiled executable (shorthand)")
	flag.BoolVar(&Debug, "debug", false, "A flag to enable debugging outputs")
	flag.BoolVar(&EmitAST, "emit-ast", false, "A flag to cause the compiler to stop at the AST generation phase and dump out the AST in a JSON file.")

	flag.Parse()
}
