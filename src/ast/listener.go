package ast

import (
	"fmt"
	"strconv"

	"github.com/noplang/nopc-go/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

//NopListener : An implementation of Antlr 4's listener interface with the aim of converting the parse tree from Antlr (which calls the exported listening functions) to the abstract syntax tree that the compiler is designed around.
type NopListener struct {
	*parser.BasenopListener

	stack stack
}

//NewNopListener : A function that creates an instance of NopListener with everything initialized to prevent null pointers.
func NewNopListener() *NopListener {
	return &NopListener{
		BasenopListener: &parser.BasenopListener{},
		stack:           make([]interface{}, 0),
	}
}

//Walk : A wrapper that takes in the Antlr parse tree and request Antlr to walk through the tree and call the listener.
//This function contains a simple panic "recovery" module (which prints out some helpful debug information and rethrow the panic).
func (n *NopListener) Walk(tree antlr.ParseTree) NopFile {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("=============")
			fmt.Println("NopListener Post-Mortem: ")
			fmt.Println("Stack:", n.stack)
			fmt.Println("Last Pop:", lastPop)
			fmt.Println("=============")
			panic(r)
		}
	}()
	walker := antlr.NewParseTreeWalker()
	walker.Walk(n, tree)
	return n.stack.Pop().(NopFile)
}

//ExitNop_file : A function that is intended to be called by Antlr when it finishes parsing the file.
//(Nop_file is the outermost node)
//Functionality: This function pops every single top-level declaration and put it into the AST properly.
func (n *NopListener) ExitNop_file(ctx *parser.Nop_fileContext) {
	nopFile := NopFile{}

	count := len(ctx.AllLet()) + len(ctx.AllConstant()) + len(ctx.AllFunction()) + len(ctx.AllStructDecl()) + len(ctx.AllImpl()) + len(ctx.AllTrait())
	declarations := n.stack.PopN(count)
	for i := 0; i < count; i++ {
		switch val := declarations[i].(type) {
		case TraitDecl:
			nopFile.Traits = append(nopFile.Traits, val)
		case ImplDecl:
			nopFile.Impls = append(nopFile.Impls, val)
		case StructDecl:
			nopFile.Structs = append(nopFile.Structs, val)
		case FunctionDecl:
			nopFile.Functions = append(nopFile.Functions, val)
		case VariableDecl:
			if val.IsConstant {
				nopFile.Constants = append(nopFile.Constants, val)
			} else {
				nopFile.Globals = append(nopFile.Globals, val)
			}
		default:
			panic(fmt.Sprintf("unexpected type: %T", declarations[i]))
		}
	}

	if ctx.ImportStatement() != nil {
		nopFile.Import = n.stack.Pop().([]string)
	}

	pragmaCount := len(ctx.AllCompiler_pragma())
	nopFile.CompilerPragma = make([]string, pragmaCount)
	poppedPragmas := n.stack.PopN(pragmaCount)
	for i, v := range poppedPragmas {
		nopFile.CompilerPragma[i] = v.(string)
	}

	n.stack.Push(nopFile)
}

//ExitCompiler_pragma : A function that is intended to be called by Antlr when it encounters a compiler pragma.
//These pragmas start with the character `#` and tells the compiler what to do with the file.
//An example of pragma is `#no_std` which disables compilation of standard library with the file.
func (n *NopListener) ExitCompiler_pragma(ctx *parser.Compiler_pragmaContext) {
	n.stack.Push(ctx.PRAGMA().GetText())
}

//ExitImportStatement :  A function that is intended to be called by Antlr when it encounters an import statement.
//It's worth noting that there could only be ONE import statement in any given file (formatted import or a single liner).
//An example is `import "std"`.
func (n *NopListener) ExitImportStatement(ctx *parser.ImportStatementContext) {
	statements := ctx.AllIDENTIFIER()
	statementCount := len(statements)
	stringStatements := make([]string, statementCount)
	for i, v := range statements {
		stringStatements[i] = v.GetText()
	}
	n.stack.Push(stringStatements)
}

type letVariableDecl struct {
	Public  bool
	Mutable bool
	Static  bool
	Name    string
}

//ExitLet :  A function that is intended to be called by Antlr when it encounters a let statement.
//This function pops out all the let variable declaration (struct name letVariableDecl) and their associated types.
//It then tries to figure out as much as possible about its type from syntax and pushes a let statement into the stack
//where it can be popped by either the global scope (ExitNop_file) or a function scope (let operation).
func (n *NopListener) ExitLet(ctx *parser.LetContext) {
	statement := n.stack.Pop().(Statement)
	letVariableCount := len(ctx.AllLetVariableDefinition())
	typeIdentifierCount := len(ctx.AllTypeIdentifier())
	count := letVariableCount + typeIdentifierCount

	currentType := VariableType{
		ToInferType: true,
	}
	listOfTargets := make([]Variable, 0, letVariableCount)
	for i := 0; i < count; i++ {
		poppedItem := n.stack.Pop()
		switch val := poppedItem.(type) {
		case letVariableDecl:
			listOfTargets = append(listOfTargets, Variable{
				Name: val.Name,
				Type: VariableType{
					Name:        currentType.Name,
					ToInferType: currentType.ToInferType,
					InnerType:   currentType.InnerType,
					Public:      val.Public,
					Mutable:     val.Mutable,
					Static:      val.Static,
				},
			})
		case VariableType:
			currentType = val
		default:
			panic(fmt.Sprintf("unexpected type: %T", poppedItem))
		}
	}
	//Reverses the array. Stack is first in last out.
	for i := 0; i < len(listOfTargets)/2; i++ {
		flippedID := len(listOfTargets) - i - 1
		listOfTargets[i], listOfTargets[flippedID] = listOfTargets[flippedID], listOfTargets[i]
	}

	n.stack.Push(VariableDecl{
		Targets:      listOfTargets,
		DefaultValue: statement,
	})
}

//ExitLetVariableDefinition :  A function that is intended to be called by Antlr when it encounters a let variable definition.
//This function reads the declaration of a variable (`pub mut someVariable`) and pushes its content onto the stack where it is
//consumed by ExitLet().
//The incentive for creating two functions instead of putting everything in ExitLet() is that Antlr's API doesn't allow inferring
//where the token is located. It is therefore needed to be separated into multiple functions to figure out if `pub` is applied to
//which variable.
func (n *NopListener) ExitLetVariableDefinition(ctx *parser.LetVariableDefinitionContext) {
	n.stack.Push(letVariableDecl{
		Public:  ctx.PUB() != nil,
		Mutable: ctx.MUT() != nil,
		Static:  ctx.STATIC() != nil,
		Name:    ctx.IDENTIFIER().GetText(),
	})
}

//ExitConstant : A function that is intended to be called by Antlr when it encounters a constant declaration.
//These are variable declarations that are only initialized once (maybe lazily if the compiler deems appropriate) and cannot be
//changed. This allows a compile time assertion such that it will never be modified.
//It is a compile error to attempt to assign something to a constant variable.
func (n *NopListener) ExitConstant(ctx *parser.ConstantContext) {
	id := ctx.AllIDENTIFIER()
	targets := make([]Variable, len(id))
	for i, v := range id {
		targets[i] = Variable{
			Name: v.GetText(),
			Type: VariableType{
				ToInferType: true,
			},
		}
	}

	defaultStatement := n.stack.Pop().(Statement)

	n.stack.Push(VariableDecl{
		Targets:      targets,
		DefaultValue: defaultStatement,
		IsConstant:   true,
	})
}

//ExitFunctionHeader : A function that is intended to be called by Antlr when it encounters a function header.
//This could happen when it's a top-level function declaration, an implementation or as a trait.
//All function parameters and return value are typed here to allow for type-checking in the next phase of compilation.
func (n *NopListener) ExitFunctionHeader(ctx *parser.FunctionHeaderContext) {
	var functionReturnType VariableType
	if ctx.ConstantType() != nil {
		functionReturnType = n.stack.Pop().(VariableType)
	}

	functionParameters := make([]FunctionParameter, 0)
	for i := 0; i < len(ctx.AllFunctionParameter()); i++ {
		count := n.stack.Pop().(int)
		for j := 0; j < count; j++ {
			poppedParam := n.stack.Pop()
			functionParameters = append(functionParameters, poppedParam.(FunctionParameter))
		}
	}

	//Reverses the array. Stack is first in last out.
	for i := 0; i < len(functionParameters)/2; i++ {
		flippedID := len(functionParameters) - i - 1
		functionParameters[i], functionParameters[flippedID] = functionParameters[flippedID], functionParameters[i]
	}

	n.stack.Push(FunctionDecl{
		Name:       ctx.IDENTIFIER().GetText(),
		Parameters: functionParameters,
		ReturnType: functionReturnType,
	})
}

//ExitFunction : A function that is intended to be called by Antlr when it encounters a function.
//This function will always be called after ExitFunctionHeader() and all the operations inside the function.
//This function simply retrieves the function header and push it back with the actual implementation (operations).
func (n *NopListener) ExitFunction(ctx *parser.FunctionContext) {
	operationCount := len(ctx.AllOperation())
	operations := make([]Operation, operationCount)
	for i, v := range n.stack.PopN(operationCount) {
		operations[i] = v.(Operation)
	}

	header := n.stack.Pop().(FunctionDecl)
	header.Body = operations
	n.stack.Push(header)
}

//ExitFunctionParameter : A function that is intended to be called by Antlr when it encounters a function parameter.
//This function is a part of ExitFunctionHeader() and is separated into multiple functions instead of aggregated with it
//to allow keyword matching to be possible as Antlr does not allow inferring token location.
func (n *NopListener) ExitFunctionParameter(ctx *parser.FunctionParameterContext) {
	if ctx.SELF() != nil {
		n.stack.Push(FunctionParameter{
			Name: "self",
			Type: VariableType{
				Name:        "self",
				ToInferType: true,
				Mutable:     ctx.MUT() != nil,
			},
		})
		n.stack.Push(1)
		return
	}

	parameterType := n.stack.Pop().(VariableType)
	identifiers := ctx.AllIDENTIFIER()
	for i := 0; i < len(identifiers); i++ {
		n.stack.Push(FunctionParameter{
			Name: identifiers[i].GetText(),
			Type: parameterType,
		})
	}
	n.stack.Push(len(identifiers))
}

//ExitOperation : A function that is intended to be called by Antlr when it encounters an operation.
//This function is called to generalize the possible operations into one Operation object which will be appended in the function
//declaration.
func (n *NopListener) ExitOperation(ctx *parser.OperationContext) {
	var operation interface{}
	var operationID OperationType
	operation = n.stack.Pop()
	switch operation.(type) {
	case VariableDecl:
		operationID = VariableDeclaration
	case ReturnOperation:
		operationID = Return
	case IfOperation:
		operationID = If
	case MatchOperation:
		operationID = Match
	case AssignmentOperation:
		operationID = Assignment
	case Statement:
		operationID = Eval
	default:
		panic(fmt.Sprintf("unexpected type: %T", operation))
	}
	n.stack.Push(Operation{
		Type:   operationID,
		Actual: operation,
	})
}

//ExitAssign : A function that is intended to be called by Antlr when it encounters an assign operation.
//This function figures out what's on the left hand side and the right hand side, apply the appropriate modifier (`+=` for example)
//and pushes it as an operation into the stack.
func (n *NopListener) ExitAssign(ctx *parser.AssignContext) {
	statement := n.stack.Pop().(Statement)
	target := n.stack.Pop().(Statement)
	operation := none
	if ctx.PLUS() != nil {
		operation = Addition
	} else if ctx.MINUS() != nil {
		operation = Subtraction
	} else if ctx.STAR() != nil {
		operation = Multiplication
	} else if ctx.SLASH() != nil {
		operation = Division
	} else if ctx.PERCENT() != nil {
		operation = Modulus
	}
	if operation != none {
		statement = Statement{
			Type: operation,
			ResolvedType: VariableType{
				ToInferType: true,
			},
			Arguments: []Statement{target, statement},
		}
	}
	n.stack.Push(AssignmentOperation{
		Target: target,
		Source: statement,
	})
}

//ExitReturnStatement : A function that is intended to be called by Antlr when it encounters a return statement.
//This function simply wraps the existing statement in an return operation so the AST shows that it is an actual return result
//and not just a evaluated and thrown away value (Like function call).
func (n *NopListener) ExitReturnStatement(ctx *parser.ReturnStatementContext) {
	n.stack.Push(ReturnOperation(n.stack.Pop().(Statement)))
}

//ExitIfStatement : A function that is intended to be called by Antlr when it encounters an if statement.
//This function wraps the if block and else block into a list of operations (similar to like how functions are a list of
//operation).
//These wrapped blocks are placed together with the condition to form the actual if statement.
func (n *NopListener) ExitIfStatement(ctx *parser.IfStatementContext) {
	elseBlock := make([]Operation, 0)
	if ctx.ElseStatement() != nil {
		elseBlock = n.stack.Pop().([]Operation)
	}

	ifBlockOperationCount := len(ctx.AllOperation())
	ifBlock := make([]Operation, ifBlockOperationCount)
	for i, v := range n.stack.PopN(ifBlockOperationCount) {
		ifBlock[i] = v.(Operation)
	}
	condition := n.stack.Pop().(Statement)

	n.stack.Push(IfOperation{
		Condition: condition,
		IfBlock:   ifBlock,
		ElseBlock: elseBlock,
	})
}

//ExitElseStatement : A function that is intended to be called by Antlr when it encounters an else statement.
//This function aggregates the operations in an else block to a list of operations. This is used to differentiate operations
//in the if block and the else block (as Antlr still doesn't allow us to know where the tokens are located).
func (n *NopListener) ExitElseStatement(ctx *parser.ElseStatementContext) {
	elseBlockOperationCount := len(ctx.AllOperation())
	elseBlock := make([]Operation, elseBlockOperationCount)
	for i, v := range n.stack.PopN(elseBlockOperationCount) {
		elseBlock[i] = v.(Operation)
	}
	n.stack.Push(elseBlock)
}

//ExitMatchStatement : A function that is intended to be called by Antlr when it encounters a match statement.
//This function aggregates a list of match candidates. It also has the condition that is actually used to match against these
//candidates.
func (n *NopListener) ExitMatchStatement(ctx *parser.MatchStatementContext) {
	matchCandidateCount := len(ctx.AllMatchCandidate())
	matchCandidates := make([]MatchCandidate, matchCandidateCount)
	for i, v := range n.stack.PopN(matchCandidateCount) {
		matchCandidates[i] = v.(MatchCandidate)
	}
	condition := n.stack.Pop().(Statement)
	n.stack.Push(MatchOperation{
		Condition: condition,
		Candidate: matchCandidates,
	})
}

//ExitMatchCandidate : A function that is intended to be called by Antlr when it encounters a match candidate.
//(A branch inside match)
//This aggregates the list of operations that should be ran if the candidate matches with the condition and the candidate to
//match.
func (n *NopListener) ExitMatchCandidate(ctx *parser.MatchCandidateContext) {
	operationCount := len(ctx.AllOperation())
	operations := make([]Operation, operationCount)
	for i, v := range n.stack.PopN(operationCount) {
		operations[i] = v.(Operation)
	}
	condition := n.stack.Pop().(Statement)
	n.stack.Push(MatchCandidate{
		Condition: condition,
		RunBlock:  operations,
	})
}

//ExitStatement : A function that is intended to be called by Antlr when it encounters a statement.
//A statement is anything that can be evaluated to a value (even if it's void).
//This function therefore contains logic to handle every single type of tokens that could be valid and holds a LHS and RHS.
//Some types of statements are handled in their respective function due to contextual difference.
func (n *NopListener) ExitStatement(ctx *parser.StatementContext) {
	operationType := none
	switch {
	case ctx.LPAREN() != nil:
		//We don't need to do anything if it's just a bracketed statement
		//Function calls are handled in ExitFunctionParameters()
		return
	case ctx.Literal() != nil:
		return //These are all handled in their own code blocks so they're no-op in this block
	case ctx.ArrayInitializer() != nil:
		return
	case ctx.ObjectInitializer() != nil:
		return
	case ctx.SELF() != nil:
		n.stack.Push(Statement{
			Type: Self,
			ResolvedType: VariableType{
				ToInferType: true,
			},
		})
		return
	case ctx.IDENTIFIER() != nil:
		n.stack.Push(Statement{
			Type: PropertyName,
			ResolvedType: VariableType{
				ToInferType: true,
			},
			LiteralValue: ctx.IDENTIFIER().GetText(),
		})
		return
	case ctx.EQUAL(0) != nil:
		//Possibilities: ==, !=, >=, <=
		switch {
		case ctx.EXCLAMATION() != nil:
			operationType = NotEqual
		case ctx.GREATERTHAN() != nil:
			operationType = GreaterEqualThan
		case ctx.LESSERTHAN() != nil:
			operationType = LesserEqualThan
		default:
			operationType = Equal
		}
	case ctx.GREATERTHAN() != nil:
		operationType = GreaterThan
	case ctx.LESSERTHAN() != nil:
		operationType = LesserThan
	case ctx.STAR() != nil:
		if len(ctx.AllStatement()) == 1 {
			pointer := n.stack.Pop().(Statement)
			n.stack.Push(Statement{
				Type: Dereference,
				ResolvedType: VariableType{
					ToInferType: true,
				},
				Arguments: []Statement{pointer},
			})
			return
		}
		operationType = Multiplication
	case ctx.SLASH() != nil:
		operationType = Division
	case ctx.PLUS() != nil:
		operationType = Addition
	case ctx.MINUS() != nil:
		operationType = Subtraction
	case ctx.PERCENT() != nil:
		operationType = Modulus
	case ctx.LBRACKET() != nil:
		operationType = ArrayAccessor
	case ctx.DOT() != nil:
		operationType = Property
	case ctx.RARROW() != nil:
		operationType = PointerAccess
	default:
		panic("exhausted all possibilities")
	}
	right := n.stack.Pop().(Statement)
	left := n.stack.Pop().(Statement)
	n.stack.Push(Statement{
		Type: operationType,
		ResolvedType: VariableType{
			ToInferType: true,
		},
		Arguments: []Statement{left, right},
	})
}

//ExitFunctionParameters : A function that is intended to be called by Antlr when it encounters function parameters.
//This function is separated from statement to differentiate between a function call without any arguments and a bracketed statement.
func (n *NopListener) ExitFunctionParameters(ctx *parser.FunctionParametersContext) {
	statementCount := len(ctx.AllStatement())
	arguments := make([]Statement, statementCount+1)
	for i, v := range n.stack.PopN(statementCount) {
		arguments[i+1] = v.(Statement)
	}
	functionToCall := n.stack.Pop().(Statement)
	arguments[0] = functionToCall
	n.stack.Push(Statement{
		Type: FunctionCall,
		ResolvedType: VariableType{
			ToInferType: true,
		},
		Arguments: arguments,
	})
}

//ExitLiteral : A function that is intended to be called by Antlr when it encounters a literal.
//This function reads the literal string and converts it into a statement that holds the literal value.
func (n *NopListener) ExitLiteral(ctx *parser.LiteralContext) {
	var value interface{}
	var typeName string
	if ctx.STRING() != nil {
		//TODO: Handle escaping
		stringValue := ctx.STRING().GetText()
		value = stringValue[1 : len(stringValue)-1] //Remove front and back quotes
		typeName = "ustring"                        //UTF-8 string
	}
	if ctx.FLOAT() != nil {
		floatValue, err := strconv.ParseFloat(ctx.FLOAT().GetText(), 64)
		if err != nil {
			panic("error reading float value: " + err.Error())
		}
		value = floatValue
		typeName = "f64"
	}
	if ctx.NUMBER() != nil {
		intValue, err := strconv.ParseInt(ctx.NUMBER().GetText(), 10, 64)
		if err != nil {
			panic("error reading int value: " + err.Error())
		}
		value = intValue
		typeName = "i64"
	}
	if ctx.TRUE() != nil {
		value = true
		typeName = "bool"
	}
	if ctx.FALSE() != nil {
		value = false
		typeName = "bool"
	}

	n.stack.Push(Statement{
		Type: Literal,
		ResolvedType: VariableType{
			Name: typeName,
		},
		LiteralValue: value,
	})
}

//ExitArrayInitializer : A function that is intended to be called by Antlr when it encounters an array initializer.
//This function reads for the array size and their default initializer.
func (n *NopListener) ExitArrayInitializer(ctx *parser.ArrayInitializerContext) {
	var arraySize int
	if ctx.NUMBER() != nil {
		var err error
		arraySize, err = strconv.Atoi(ctx.NUMBER().GetText())
		if err != nil {
			panic("error reading array initializer size: " + err.Error())
		}
	} else {
		arraySize = len(ctx.AllStatement())
	}
	//First is the type of array to create
	//Second is the length of the array to create
	arguments := make([]Statement, len(ctx.AllStatement())+2)
	arguments[1] = Statement{
		Type:         Literal,
		LiteralValue: arraySize,
		ResolvedType: VariableType{
			Name: "i64",
		},
	}
	for i, v := range n.stack.PopN(len(ctx.AllStatement())) {
		arguments[i+2] = v.(Statement)
	}
	arguments[0] = Statement{
		Type:         Type,
		ResolvedType: n.stack.Pop().(VariableType),
	}
	n.stack.Push(Statement{
		Type: ArrayInitializer,
		ResolvedType: VariableType{
			ToInferType: true,
		},
		Arguments: arguments,
	})
}

//ExitObjectInitializer : A function that is intended to be called by Antlr when it encounters an object initializer.
//This function reads for the type of object to initialize and its default property values.
func (n *NopListener) ExitObjectInitializer(ctx *parser.ObjectInitializerContext) {
	arguments := make([]Statement, len(ctx.AllStatement())+1)
	for i, v := range n.stack.PopN(len(ctx.AllStatement())) {
		arguments[i+1] = v.(Statement) //First is the type of object to create
	}
	objectType := n.stack.Pop().(VariableType)
	arguments[0] = Statement{
		Type:         Type,
		ResolvedType: objectType,
	}
	n.stack.Push(Statement{
		Type: ObjectInitializer,
		//We already know that object type is being initialized. No point in making the compiler infer it.
		ResolvedType: objectType,
		Arguments:    arguments,
	})
}

//ExitStructDecl : A function that is intended to be called by Antlr when it encounters a struct declaration.
//It simply searches for a list of fields in the struct so typing could be done properly.
func (n *NopListener) ExitStructDecl(ctx *parser.StructDeclContext) {
	//Cap is set at the amount of struct var. There's at least that much anyway.
	fieldNames := make([]string, 0, len(ctx.AllStructVar()))
	fieldTypes := make([]VariableType, 0, len(ctx.AllStructVar()))
	for i := 0; i < len(ctx.AllStructVar()); i++ {
		count := n.stack.Pop().(int)
		poppedNames := make([]string, count)
		poppedTypes := make([]VariableType, count)
		for j, v := range n.stack.PopN(count) {
			field := v.(Variable)
			poppedNames[j] = field.Name
			poppedTypes[j] = field.Type
		}
		fieldNames = append(poppedNames, fieldNames...)
		fieldTypes = append(poppedTypes, fieldTypes...)
	}
	structName := ctx.IDENTIFIER().GetText()
	n.stack.Push(StructDecl{
		Public:     ctx.PUB() != nil,
		Name:       structName,
		FieldNames: fieldNames,
		FieldTypes: fieldTypes,
	})
}

//ExitStructVar : A function that is intended to be called by Antlr when it encounters a struct variable.
//This is a supporting function to ExitStructDecl() to make multiple identifiers having the same type possible.
//(Antlr still doesn't allow inferring token location)
func (n *NopListener) ExitStructVar(ctx *parser.StructVarContext) {
	fieldType := n.stack.Pop().(VariableType)
	if ctx.PUB() != nil {
		fieldType.Public = true
	}
	identifiers := ctx.AllIDENTIFIER()
	for _, v := range identifiers {
		n.stack.Push(Variable{
			Name: v.GetText(),
			Type: fieldType,
		})
	}
	n.stack.Push(len(identifiers))
}

//ExitImpl : A function that is intended to be called by Antlr when it encounters an implementation definition.
//This removes functions from the top-level context and place it inside of a implementation struct to show that it's part of
//an implementation only.
func (n *NopListener) ExitImpl(ctx *parser.ImplContext) {
	implCount := len(ctx.AllFunction())
	functions := make([]FunctionDecl, implCount)
	for i, v := range n.stack.PopN(implCount) {
		functions[i] = v.(FunctionDecl)
	}
	var implTypeTarget VariableType
	if ctx.TypeIdentifier() != nil {
		implTypeTarget = n.stack.Pop().(VariableType)
	}
	name := ctx.IDENTIFIER().GetText()
	n.stack.Push(ImplDecl{
		TraitName:       name,
		TypeName:        implTypeTarget,
		Implementations: functions,
	})
}

//ExitTrait : A function that is intended to be called by Antlr when it encounters a trait declaration.
//This function reads function headers that has been discovered and groups it under itself to show that it's a trait declaration
//and not just an empty function.
func (n *NopListener) ExitTrait(ctx *parser.TraitContext) {
	headerCount := len(ctx.AllFunctionHeader())
	functionHeaders := make([]FunctionDecl, headerCount)
	for i, v := range n.stack.PopN(headerCount) {
		functionHeaders[i] = v.(FunctionDecl)
	}
	traitName := ctx.IDENTIFIER().GetText()
	n.stack.Push(TraitDecl{
		Name:      traitName,
		Functions: functionHeaders,
	})
}

//ExitTypeIdentifier : A function that is intended to be called by Antlr when it encounters a (potentially mutable) type identifier.
//This function simply sets mutable to true if it encounters a mutable type.
//The main type processing is done in ExitConstantType().
func (n *NopListener) ExitTypeIdentifier(ctx *parser.TypeIdentifierContext) {
	if ctx.MUT() == nil {
		return //Do nothing. It's already the correct type.
	}
	toMutateType := n.stack.Pop().(VariableType)
	toMutateType.Mutable = true
	n.stack.Push(toMutateType)
}

//ExitConstantType : A function that is intended to be called by Antlr when it encounters a non-mutable type identifier.
//This function figures out the type of an object by recursively adding onto it.
//Types are set to be mutable in ExitTypeIdentifier().
//They are separated into two functions to enforce mutability rules in the grammar file.
func (n *NopListener) ExitConstantType(ctx *parser.ConstantTypeContext) {
	if ctx.STAR() != nil {
		pointedToType := n.stack.Pop().(VariableType)
		n.stack.Push(VariableType{
			Name:      "Pointer",
			InnerType: &pointedToType,
		})
		return
	}
	if ctx.LBRACKET() != nil {
		contentType := n.stack.Pop().(VariableType)
		n.stack.Push(VariableType{
			Name:      "Array",
			InnerType: &contentType,
		})
		return
	}
	var innerType *VariableType
	if ctx.LESSERTHAN() != nil {
		castedType := n.stack.Pop().(VariableType)
		innerType = &castedType
	}
	n.stack.Push(VariableType{
		Name:      ctx.IDENTIFIER().GetText(),
		InnerType: innerType,
	})
}
