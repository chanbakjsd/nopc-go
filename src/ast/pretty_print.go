package ast

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
)

//PrettyPrintAST : A function that prints a given parsed AST and serializes it into somewhat readable JSON format.
func PrettyPrintAST(n NopFile) string {
	s, err := json.MarshalIndent(n, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(s)
}

//PrettyPrintASTToFile : A utility function that calls PrettyPrintAST() and saves the resulting string in a file.
func PrettyPrintASTToFile(n NopFile) {
	result := PrettyPrintAST(n)
	file, err := os.Create(n.Source + ".ast")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(file)
	w.WriteString(result)
	w.Flush()
}

//String : A function that allows operation type to be displayed in its actual enum name instead of the underlying value.
func (o OperationType) String() string {
	switch o {
	case VariableDeclaration:
		return "[VariableDecl]"
	case Return:
		return "[ReturnStatement]"
	case If:
		return "[IfStatement]"
	case Match:
		return "[MatchStatement]"
	case Assignment:
		return "[Assignment]"
	case Eval:
		return "[EvaluatedStatement]"
	default:
		return "[OperationType " + strconv.Itoa(int(o)) + "]"
	}
}

//MarshalJSON : A function that takes the enum name of an operation type and turns it into the format Go's JSON library need.
func (o OperationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

//String : A function that allows statement types to be displayed in its actual enum name instead of the underlying value.
func (s StatementType) String() string {
	switch s {
	case Equal:
		return "=="
	case NotEqual:
		return "!="
	case LesserThan:
		return "<"
	case LesserEqualThan:
		return "<="
	case GreaterThan:
		return ">"
	case GreaterEqualThan:
		return ">="
	case FunctionCall:
		return "[FunctionCall]"
	case Multiplication:
		return "[Multiplication]"
	case Division:
		return "[Division]"
	case Addition:
		return "[Addition]"
	case Subtraction:
		return "[Subtraction]"
	case Modulus:
		return "[Modulus]"
	case Self:
		return "[Self]"
	case Literal:
		return "[Literal]"
	case ArrayInitializer:
		return "[ArrayInitializer]"
	case ObjectInitializer:
		return "[ObjectInitializer]"
	case ArrayAccessor:
		return "[ArrayAccessor]"
	case PointerAccess:
		return "[PointerAccess]"
	case Dereference:
		return "[Dereference]"
	case Property:
		return "[Property]"
	case PropertyName:
		return "[PropertyName]"
	case Type:
		return "[Type]"
	case none:
		return "[INTERNAL none]"
	default:
		return "[StatementType " + strconv.Itoa(int(s)) + "]"
	}
}

//MarshalJSON : A function that takes the enum name of a statement type and turns it into the format Go's JSON library need.
func (s StatementType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
