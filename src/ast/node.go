package ast

type OperationType int8
type StatementType int8

const (
	VariableDeclaration OperationType = iota
	Return
	If
	Match
	Assignment
	Eval
)

const (
	Equal StatementType = iota
	NotEqual
	LesserThan
	LesserEqualThan
	GreaterThan
	GreaterEqualThan
	FunctionCall
	Multiplication
	Division
	Addition
	Subtraction
	Modulus
	Self
	Literal
	ArrayInitializer
	ObjectInitializer
	ArrayAccessor
	PointerAccess
	Dereference
	Property
	PropertyName
	Type
	none //A placeholder used in code to denote that the value is used as-is
)

type NopFile struct {
	Source         string
	CompilerPragma []string
	Import         []string
	Structs        []StructDecl
	Traits         []TraitDecl
	Impls          []ImplDecl
	Globals        []VariableDecl
	Constants      []VariableDecl //Globals but cannot be modified in runtime
	Functions      []FunctionDecl
}

type StructDecl struct {
	Public     bool
	Name       string
	FieldNames []string
	FieldTypes []VariableType
}

type TraitDecl struct {
	Name      string
	Functions []FunctionDecl
}

type ImplDecl struct {
	TraitName       string
	TypeName        VariableType
	Implementations []FunctionDecl
}

type VariableDecl struct {
	Targets      []Variable
	DefaultValue Statement
	IsConstant   bool
}

type Variable struct {
	Name string
	Type VariableType
}

type VariableType struct {
	Name        string
	ToInferType bool          //Requests the compiler to infer Name and InnerType
	InnerType   *VariableType //Outer<InnerType>. Pointer and Arrays go here too.
	Public      bool
	Static      bool
	Mutable     bool
}

type FunctionDecl struct {
	Name       string
	Parameters []FunctionParameter
	ReturnType VariableType
	Body       []Operation
}

type FunctionParameter struct {
	Name string
	Type VariableType
}

type Operation struct {
	Type   OperationType
	Actual interface{}
}

type ReturnOperation Statement

type IfOperation struct {
	Condition Statement
	IfBlock   []Operation
	ElseBlock []Operation
}

type MatchOperation struct {
	Condition Statement
	Candidate []MatchCandidate
}

type MatchCandidate struct {
	Condition Statement
	RunBlock  []Operation
}

type AssignmentOperation struct {
	Target Statement
	Source Statement
}

type Statement struct {
	Type         StatementType
	ResolvedType VariableType
	LiteralValue interface{}
	Arguments    []Statement
}
