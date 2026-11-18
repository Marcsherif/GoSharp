package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "."

	IDENT  = "IDENT"
	INT    = "INT"
	DOUBLE = "DOUBLE"

	PERCENT  = "%"
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	LE = "<="
	LT = "<"
	GT = ">"
	GE = ">="

	EQ     = "=="
	NOT_EQ = "<>"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "["
	RBRACE = "]"

	KEYWORD = "KEYWORD"
	AND     = "AND"
	DEF     = "DEF"
	DO      = "DO"
	ELSE    = "ELSE"
	FED     = "FED"
	FI      = "FI"
	IF      = "IF"
	NOT     = "NOT"
	OD      = "OD"
	OR      = "OR"
	PRINT   = "PRINT"
	RETURN  = "RETURN"
	THEN    = "THEN"
)

var keywords = map[string]struct{}{
	"and":		{},
	"def":		{},
	"do":		{},
	"double":	{},
	"else":		{},
	"fed":		{},
	"fi":		{},
	"if":		{},
	"int":		{},
	"not":		{},
	"od":		{},
	"or":		{},
	"print":	{},
	"return":	{},
	"then":		{},
	"while":	{},
}

type Token struct {
	Type    TokenType
	Literal string
}

func LookupIdent(ident string) TokenType {
	if _, ok := keywords[ident]; ok {
		return KEYWORD
	}
	return IDENT
}
