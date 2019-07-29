package ast

import (
	"fmt"
	"strconv"

	"github.com/noplang/nopc-go/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type NopListener struct {
	*parser.BasenopListener

	stack stack
}

func NewNopListener() *NopListener {
	return &NopListener{
		BasenopListener: &parser.BasenopListener{},
		stack:           make([]interface{}, 0),
	}
}

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

func (n *NopListener) ExitCompiler_pragma(ctx *parser.Compiler_pragmaContext) {
	n.stack.Push(ctx.PRAGMA().GetText())
}

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

func (n *NopListener) ExitLetVariableDefinition(ctx *parser.LetVariableDefinitionContext) {
	n.stack.Push(letVariableDecl{
		Public:  ctx.PUB() != nil,
		Mutable: ctx.MUT() != nil,
		Static:  ctx.STATIC() != nil,
		Name:    ctx.IDENTIFIER().GetText(),
	})
}

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

func (n *NopListener) ExitFunctionHeader(ctx *parser.FunctionHeaderContext) {
	var functionReturnType VariableType
	if ctx.ConstantType() != nil {
		functionReturnType = n.stack.Pop().(VariableType)
	}

	functionParameters := make([]FunctionParameter, 0)
	for i := 0; i < len(ctx.AllFunctionParameter()); i++ {
		count := n.stack.Pop().(int)
		poppedParam := n.stack.PopN(count)
		for j := 0; j < count; j++ {
			functionParameters = append(functionParameters, poppedParam[j].(FunctionParameter))
		}
	}
	n.stack.Push(FunctionDecl{
		Name:       ctx.IDENTIFIER().GetText(),
		Parameters: functionParameters,
		ReturnType: functionReturnType,
	})
}

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

func (n *NopListener) ExitReturnStatement(ctx *parser.ReturnStatementContext) {
	n.stack.Push(ReturnOperation(n.stack.Pop().(Statement)))
}

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

func (n *NopListener) ExitElseStatement(ctx *parser.ElseStatementContext) {
	elseBlockOperationCount := len(ctx.AllOperation())
	elseBlock := make([]Operation, elseBlockOperationCount)
	for i, v := range n.stack.PopN(elseBlockOperationCount) {
		elseBlock[i] = v.(Operation)
	}
	n.stack.Push(elseBlock)
}

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

func (n *NopListener) ExitStatement(ctx *parser.StatementContext) {
	operationType := none
	switch {
	case ctx.LPAREN() != nil:
		statementCount := len(ctx.AllStatement())
		if statementCount != 1 {
			//We don't need to do anything if it's just a bracketed statement
			arguments := make([]Statement, statementCount)
			for i, v := range n.stack.PopN(statementCount) {
				arguments[i] = v.(Statement)
			}
			n.stack.Push(Statement{
				Type:      FunctionCall,
				Arguments: arguments,
			})
		}
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
				Type:      Dereference,
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
		Type:      ArrayInitializer,
		Arguments: arguments,
	})
}

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
		Type:      ObjectInitializer,
		Arguments: arguments,
	})
}

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

func (n *NopListener) ExitStructVar(ctx *parser.StructVarContext) {
	fieldType := n.stack.Pop().(VariableType)
	identifiers := ctx.AllIDENTIFIER()
	for _, v := range identifiers {
		n.stack.Push(Variable{
			Name: v.GetText(),
			Type: fieldType,
		})
	}
	n.stack.Push(len(identifiers))
}

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

func (n *NopListener) ExitTypeIdentifier(ctx *parser.TypeIdentifierContext) {
	if ctx.MUT() == nil {
		return //Do nothing. It's already the correct type.
	}
	toMutateType := n.stack.Pop().(VariableType)
	toMutateType.Mutable = true
	n.stack.Push(toMutateType)
}

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
