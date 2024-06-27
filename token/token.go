package token

import (
	"bufio"
	"os"
)

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

type Token struct {
	Type    TokenType
	Literal string
}

func LookupIdent(ident string) TokenType {
	file, err := os.Open("./keywords/keywords.txt")
    _ = err
	defer file.Close()

	var keywords []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		keyword := scanner.Text()
		keywords = append(keywords, keyword)
	}

    for _, keyword := range keywords {
        if ident == keyword {
            return KEYWORD
        }
    }
	return IDENT
}
