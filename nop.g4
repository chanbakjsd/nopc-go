grammar nop;

nop_file: compiler_pragma* importStatement? (let | constant | function | structDecl | impl | trait)+;

compiler_pragma: PRAGMA;

importStatement
: IMPORT LPAREN (IDENTIFIER COMMA)* IDENTIFIER RPAREN
| IMPORT IDENTIFIER
;

let
: LET letVariableDefinition typeIdentifier? EQUAL statement
| LET LPAREN (letVariableDefinition typeIdentifier? COMMA)+ letVariableDefinition typeIdentifier? RPAREN EQUAL statement
;

letVariableDefinition: PUB? MUT? STATIC? IDENTIFIER;

constant
: CONST PUB? IDENTIFIER EQUAL statement
| CONST LPAREN (PUB? IDENTIFIER COMMA)+ PUB? IDENTIFIER RPAREN EQUAL statement
;

functionHeader: FN IDENTIFIER LPAREN ((functionParameter COMMA)* functionParameter)? RPAREN constantType?;
function: functionHeader LCURLY operation* RCURLY;
functionParameter
: (IDENTIFIER COMMA)* IDENTIFIER typeIdentifier //x, y, z i32
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

assign: statement (PLUS | MINUS | STAR | SLASH | PERCENT)? EQUAL statement;
returnStatement: RETURN statement;
ifStatement: IF statement LCURLY operation* RCURLY elseStatement?;
elseStatement: ELSE LCURLY operation* RCURLY;
matchStatement: MATCH statement LCURLY matchCandidate+ RCURLY;
matchCandidate: statement RARROW ((LCURLY operation* RCURLY) | operation);

statement
: LPAREN statement RPAREN
| literal
| IDENTIFIER
| STAR statement
| SELF
| statement (DOT | RARROW) statement
| statement LPAREN functionParameters RPAREN
| statement EQUAL EQUAL statement
| statement EXCLAMATION EQUAL statement
| statement LESSERTHAN (EQUAL)? statement
| statement GREATERTHAN (EQUAL)? statement
| statement (STAR | SLASH | PERCENT) statement
| statement (PLUS | MINUS) statement
| statement LBRACKET statement RBRACKET
| arrayInitializer
| objectInitializer
;

functionParameters: (statement (COMMA statement)*)?;

literal
: STRING
| FLOAT
| NUMBER
| TRUE
| FALSE
;

arrayInitializer
: LBRACKET NUMBER RBRACKET typeIdentifier
| LBRACKET NUMBER? RBRACKET typeIdentifier LCURLY (statement COMMA)* statement RCURLY
;

objectInitializer: constantType LCURLY ((statement COMMA)* statement)? RCURLY;

structDecl: PUB? STRUCT IDENTIFIER LCURLY (structVar COMMA)* structVar RCURLY;
structVar: PUB? (IDENTIFIER COMMA)* IDENTIFIER typeIdentifier;

impl: IMPL IDENTIFIER (FOR typeIdentifier)? LCURLY function* RCURLY;
trait: TRAIT IDENTIFIER LCURLY functionHeader* functionHeader RCURLY;

typeIdentifier
: MUT constantType
| constantType
;
constantType
: STAR typeIdentifier
| LBRACKET RBRACKET typeIdentifier
| IDENTIFIER (LESSERTHAN typeIdentifier GREATERTHAN)?
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
