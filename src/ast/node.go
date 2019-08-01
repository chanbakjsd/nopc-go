package ast

//OperationType : A constant that defines the type of operation we're dealing with to allow information to be extracted by
//converting the operation to the correct type.
type OperationType int8

//StatementType : A constant that defines the type of statement we're dealing with to allow the compiler to generate the correct
//code for the operation.
type StatementType int8

const (
	//VariableDeclaration : An operation that declares a variable. (Let)
	VariableDeclaration OperationType = iota
	//Return : An operation that returns a value (ie. returning from a function).
	Return
	//If : An operation that carries out a if/else operation.
	If
	//Match : An operation that carries out a match operation.
	Match
	//Assignment : An operation that changes the value of a property with an evaluated value.
	Assignment
	//Eval : An operation that evaluates and throw away the resulting value to use its side effect.
	//An example is a function call.
	Eval
)

const (
	//Equal : This statement compares if LHS is equal to RHS.
	Equal StatementType = iota
	//NotEqual : This statement compares if LHS is not equal to RHS.
	NotEqual
	//LesserThan : This statement compares if LHS is lesser than RHS.
	LesserThan
	//LesserEqualThan : This statement compares if LHS is lesser than or equal to RHS.
	LesserEqualThan
	//GreaterThan : This statement compares if LHS is greater than RHS.
	GreaterThan
	//GreaterEqualThan : This statement compares if LHS is greater than or equal to RHS.
	GreaterEqualThan
	//FunctionCall : This statement calls a function. Argument is stored in the format: [FunctionToCall, Arguments...]
	FunctionCall
	//Multiplication : This statement returns LHS multiplied by RHS.
	Multiplication
	//Division : This statement returns LHS divided by RHS.
	Division
	//Addition : This statement returns RHS added to LHS.
	Addition
	//Subtraction : This statement returns RHS subtracted from LHS.
	Subtraction
	//Modulus : This statement returns LHS mod RHS, keeping the sign.
	Modulus
	//Self : This is a statement without any arguments. It simply refers to oneself. Used in Impl declarations.
	Self
	//Literal : A statement that can be evaluated to a value.
	Literal
	//ArrayInitializer : A statement that initializes an array.
	ArrayInitializer
	//ObjectInitializer : A statement that initializes an object.
	ObjectInitializer
	//ArrayAccessor : A statement that accesses the value of an array. (array[i])
	ArrayAccessor
	//PointerAccess : A statement that skips the dereferencing step and reads the struct value. (a->b)
	PointerAccess
	//Dereference : A statement that takes in a pointer and return the actual object.
	Dereference
	//Property : A statement that access a property of an object (like a field in a struct).
	Property
	//PropertyName : These are placeholders to show which fields to access. They'll evaluate to a value if the field exist.
	PropertyName
	//Type : A statement that tells ArrayInitializer and ObjectInitializer what type to initialize.
	Type
	none //A placeholder used in code to denote that the value is used as-is
)

//NopFile : The parsed AST. It contains all info about a file.
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

//StructDecl : A declaration of a struct.
type StructDecl struct {
	Public     bool
	Name       string
	FieldNames []string
	FieldTypes []VariableType
}

//TraitDecl : A declaration of a trait.
type TraitDecl struct {
	Name      string
	Functions []FunctionDecl
}

//ImplDecl : A declaration of an implementation (impl).
type ImplDecl struct {
	TraitName       string
	TypeName        VariableType
	Implementations []FunctionDecl
}

//VariableDecl : A declaration of variable (let).
type VariableDecl struct {
	Targets      []Variable
	DefaultValue Statement
	IsConstant   bool
}

//Variable : A struct that holds information about a variable (that has been declared on is being assigned to).
type Variable struct {
	Name string
	Type VariableType
}

//VariableType : A struct that holds information about a type.
type VariableType struct {
	Name        string
	ToInferType bool          //Requests the compiler to infer Name and InnerType
	InnerType   *VariableType //Outer<InnerType>. Pointer and Arrays go here too.
	Public      bool
	Static      bool
	Mutable     bool
}

//FunctionDecl : A declaration of a function (fn).
type FunctionDecl struct {
	Name       string
	Parameters []FunctionParameter
	ReturnType VariableType
	Body       []Operation
}

//FunctionParameter : A struct that holds the information about a specific function parameter.
type FunctionParameter struct {
	Name string
	Type VariableType
}

//Operation : An operation that should be ran by a function.
type Operation struct {
	Type   OperationType
	Actual interface{}
}

//ReturnOperation : An operation that is used in a function to return.
type ReturnOperation Statement

//IfOperation : An operation that checks if the statement condition is true and run code conditionally.
type IfOperation struct {
	Condition Statement
	IfBlock   []Operation
	ElseBlock []Operation
}

//MatchOperation : An operation that chooses the first candidate to match with the list of candidates.
type MatchOperation struct {
	Condition Statement
	Candidate []MatchCandidate
}

//MatchCandidate : A candidate to be matched by the MatchOperation.
type MatchCandidate struct {
	Condition Statement
	RunBlock  []Operation
}

//AssignmentOperation : An operation that assigns a value (result of statement) to a field or a variable.
type AssignmentOperation struct {
	Target Statement
	Source Statement
}

//Statement : A struct containing anything that can be evaluated to a value.
type Statement struct {
	Type         StatementType
	ResolvedType VariableType
	LiteralValue interface{}
	Arguments    []Statement
}
