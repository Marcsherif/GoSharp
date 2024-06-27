package lexer

import (
	"bufio"
	"ezsharp/token"
	"os"
	"strconv"
	"strings"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
    line         int
}

func readState() ([][]int) {
    file, err := os.Open("./Table/table.txt")
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    check(err)
    n := 21
    myArray := make([][]int, n)

    i := 0
    for scanner.Scan() {
        line := scanner.Text()
        values := strings.Fields(line)
        myArray[i] = make([]int, 128)

        for j, valueStr := range values {
            value, err := strconv.Atoi(valueStr)
            check(err)
            myArray[i][j] = value
        }
        i++
        if i == n {
            break
        }
    }

    return myArray
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() (token.Token, int) {
	var tok token.Token


    TT := readState()
    l.skipWhitespace()
    state := 0
    str := ""
    for ((state >= 0) && (state <= 22)) {
        switch (state) {
        case 22:
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
            goto next
        case 0:
            state = TT[state][int(l.ch)]

            if l.ch == '%'{
                tok = newToken(token.PERCENT)
                goto next
            } else if l.ch == '*' {
                tok = newToken(token.ASTERISK)
                goto next
            } else if l.ch == '+' {
                tok = newToken(token.PLUS)
                goto next
            } else if l.ch == '-' {
                tok = newToken(token.MINUS)
                goto next
            } else if l.ch == '/' {
                tok = newToken(token.SLASH)
                goto next
            } else if l.ch == ';' {
                tok = newToken(token.SEMICOLON)
                goto next
            } else if l.ch == ',' {
                tok = newToken(token.COMMA)
                goto next
            } else if l.ch == '.' {
                tok = newToken(token.EOF)
                goto next
            } else if l.ch == '(' {
                tok = newToken(token.LPAREN)
                goto next
            } else if l.ch == ')' {
                tok = newToken(token.RPAREN)
                goto next
            } else if l.ch == '[' {
                tok = newToken(token.LBRACE)
                goto next
            } else if l.ch == ']' {
                tok = newToken(token.RBRACE)
                goto next
            }
        case 1:
            l.readChar()
            l.skipWhitespace()
            state = TT[state][int(l.ch)]
        case 2:
            tok = newToken(token.LE)
			goto next
        case 3:
            tok = newToken(token.NOT_EQ)
			goto next
        case 4:
            tok = newToken(token.LT)
            l.backChar()
            goto next
        case 5:
            l.readChar()
            l.skipWhitespace()
            state = TT[state][int(l.ch)]
        case 6:
            tok = newToken(token.EQ)
            goto next
        case 7:
            tok = newToken(token.ASSIGN)
            l.backChar()
            goto next
        case 8:
            l.readChar()
            l.skipWhitespace()
            state = TT[state][int(l.ch)]
        case 9:
            tok = newToken(token.GE)
			goto next
        case 10:
            tok = newToken(token.GT)
            l.backChar()
			goto next
        case 11:
            str = str + string([]byte{l.ch})
            l.readChar()
            state = TT[state][int(l.ch)]
        case 12:
            tok.Literal = str
            tok.Type = token.LookupIdent(tok.Literal)
            l.backChar()
            goto next
        case 13:
            str = str + string([]byte{l.ch})
            l.readChar()
            state = TT[state][int(l.ch)]
        case 14:
            str = str + string([]byte{l.ch})
            l.readChar()
            l.skipWhitespace()
            state = TT[state][int(l.ch)]
        case 15:
            str = str + string([]byte{l.ch})
            l.readChar()
            state = TT[state][int(l.ch)]
        case 16:
            str = str + string([]byte{l.ch})
            l.readChar()
            state = TT[state][int(l.ch)]
        case 17:
            str = str + string([]byte{l.ch})
            l.readChar()
            state = TT[state][int(l.ch)]
        case 18:
            str = str + string([]byte{l.ch})
            l.readChar()
            state = TT[state][int(l.ch)]
        case 19:
            tok.Literal = str
            tok.Type = token.DOUBLE
            l.backChar()
            goto next
        case 20:
            tok.Literal = str
            tok.Type = token.INT
            l.backChar()
            goto next
        case 21:
            tok.Literal = str
            tok.Type = token.DOUBLE
            l.backChar()
            goto next
        }
    }
    next:
    l.readChar()
    return tok, l.line
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        if l.ch == '\n' {
            l.line++
        }
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) backChar() {
	if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.position -= 1
        l.readPosition -= 1
        l.ch = l.input[l.position]
    }
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenType)}
}
