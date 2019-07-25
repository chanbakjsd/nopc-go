package flags

import "fmt"

//DebugPrint : Logs a message if we're being ran in debug mode.
func DebugPrint(message string) {
	if Debug {
		fmt.Println(message)
	}
}
