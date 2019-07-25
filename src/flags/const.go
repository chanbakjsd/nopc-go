package flags

import "flag"

var (
	OutputLocation string
	Debug          bool
)

func init() {
	flag.StringVar(&OutputLocation, "output", "<auto>", "Where to output the compiled executable")
	flag.StringVar(&OutputLocation, "o", "<auto>", "Where to output the compiled executable (shorthand)")
	flag.BoolVar(&Debug, "debug", false, "A flag to enable debugging outputs")

	flag.Parse()
}
