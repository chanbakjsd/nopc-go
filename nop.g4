grammar nop;

nop_file: compiler_pragma* importStatement? (let | constant | function | structDecl | impl | trait)+;

compiler_pragma: PRAGMA;

importStatement
: IMPORT LPAREN (IDENTIFIER COMMA)* IDENTIFIER RPAREN
| IMPORT IDENTIFIER
;

let
: LET letVariableDefinition type? EQUAL statement
| LET LPAREN (letVariableDefinition type? COMMA)+ letVariableDefinition type? RPAREN EQUAL statement
;

letVariableDefinition: PUB? MUT? STATIC? IDENTIFIER;

constant
: CONST PUB? IDENTIFIER EQUAL statement
| CONST LPAREN (PUB? IDENTIFIER COMMA)+ PUB? IDENTIFIER RPAREN EQUAL statement
;

functionHeader: FN IDENTIFIER LPAREN ((functionParameter COMMA)* functionParameter)? RPAREN IDENTIFIER?;
function: functionHeader LCURLY operation* RCURLY;
functionParameter
: (IDENTIFIER COMMA)* IDENTIFIER type //x, y, z i32
| STAR MUT? SELF
;

operation
: let
| returnStatement
| ifStatement
| matchStatement
| assign
| statement
;

assign: property (PLUS | MINUS | STAR | SLASH | PERCENT)? EQUAL statement;
returnStatement: RETURN statement;
ifStatement: IF statement LCURLY operation* RCURLY (ELSE LCURLY operation* RCURLY)?;
matchStatement: MATCH statement LCURLY (statement RARROW ((LCURLY operation* RCURLY) | operation))+ RCURLY;

statement
: LPAREN statement RPAREN
| statement EQUAL EQUAL statement
| statement EXCLAMATION EQUAL statement
| statement LESSERTHAN (EQUAL)? statement
| statement GREATERTHAN (EQUAL)? statement
| functionCall
| statement (STAR | SLASH | PERCENT) statement
| statement (PLUS | MINUS) statement
| SELF
| literal
| arrayInitializer
| objectInitializer
| arrayAccessor
| property
;

literal
: STRING
| FLOAT
| NUMBER
| TRUE
| FALSE
;

arrayInitializer
: LBRACKET NUMBER RBRACKET type
| LBRACKET NUMBER? RBRACKET type LCURLY (statement COMMA)* statement RCURLY
;

objectInitializer: IDENTIFIER LCURLY ((statement COMMA)* statement)? RCURLY;

arrayAccessor: property LBRACKET statement RBRACKET;
property: (STAR* SELF (DOT | RARROW))? (STAR* IDENTIFIER (DOT | RARROW))* STAR* IDENTIFIER;
functionCall: property LPAREN ((statement COMMA)* statement)? RPAREN;

structDecl: PUB? STRUCT IDENTIFIER LCURLY (structVar COMMA)* structVar RCURLY;
structVar: PUB? (IDENTIFIER COMMA)* IDENTIFIER type;

impl: IMPL IDENTIFIER (FOR IDENTIFIER)? LCURLY function* RCURLY;
trait: TRAIT IDENTIFIER LCURLY functionHeader* functionHeader RCURLY;

type
: MUT constantType
| constantType
;
constantType
: STAR type
| LBRACKET RBRACKET type
| IDENTIFIER (LESSERTHAN type GREATERTHAN)?
;

/** LEXER LOGIC */
ONELINECOMMENT: '//' .*? [\r\n] -> skip;
MULTILINECOMMENT: '/*' .*? '*/' -> skip;
STRING: ('"' (.*? ~[\\])+? '"') | ('\'' (.*? ~[\\])+? '\'');
PRAGMA: '#' ~[\r\n]*;
IMPORT: 'import';
IMPL: 'impl';
TRAIT: 'trait';
LET: 'let';
CONST: 'const';
PUB: 'pub';
MUT: 'mut';
STATIC: 'static';
FN: 'fn';
RETURN: 'return';
IF: 'if';
ELSE: 'else';
MATCH: 'match';
TRUE: 'true';
FALSE: 'false';
STRUCT: 'struct';
FOR: 'for';
SELF: 'self';
COMMA: ',';
LPAREN: '(';
RPAREN: ')';
LCURLY: '{';
RCURLY: '}';
LBRACKET: '[';
RBRACKET: ']';
EQUAL: '=';
DOT: '.';
RARROW: '->';
PLUS: '+';
MINUS: '-';
STAR: '*';
SLASH: '/';
PERCENT: '%';
EXCLAMATION: '!';
LESSERTHAN: '<';
GREATERTHAN: '>';
IDENTIFIER: [A-Za-z][A-Za-z0-9]*;
ALPHA: [A-Za-z];
NUMBER: [0-9]+;
FLOAT: NUMBER DOT NUMBER;
DIGIT: [0-9];
WHITESPACE: [ \t\r\n]+ -> skip; //Whitespace doesn't hold any meaning in our syntax
