package main

import "github.com/antlr/antlr4/runtime/Go/antlr"

//errorListener : An error listener that records if any error has been encountered.
type errorListener struct {
	EncounteredError bool
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.EncounteredError = true
}
func (l *errorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	//No-op. This function is here to implement the interface.
}
func (l *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	//No-op. This function is here to implement the interface.
}
func (l *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	//No-op. This function is here to implement the interface.
}
