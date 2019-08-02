package ast

import (
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func TestSyntaxError(t *testing.T) {
	//Simply invoke syntax error to be sure it actually works.
	//More of a sanity check to make sure this isn't the root cause of any future problem.

	listener := errorListener{}
	if listener.EncounteredError {
		t.Error("Initialized error listener has unexpectedly already listened to an error occurring.")
	}

	//Faking is hard.
	recognizer := antlr.BaseLexer{} //BaseRecognizer doesn't actually implement Recognizer. You need Lexer or Parser.
	recognitionException := antlr.BaseRecognitionException{}
	listener.SyntaxError(&recognizer, "", 1, 2, "", &recognitionException)

	if !listener.EncounteredError {
		t.Error("Error listener should have the error encountered flag set!")
	}

	noDuplicateSanity := errorListener{}
	if noDuplicateSanity.EncounteredError {
		//This could really only happen if the ErrorEncountered flag is shared which definitely shouldn't happen.
		t.Error("Initialized error listener has unexpectedly already listened to an error occurring.")
	}
}
