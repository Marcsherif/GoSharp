package lexer

import (
	"bufio"
	"ezsharp/token"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var TT [][]int

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
    line         int
}

func loadTransitionTable() () {
    file, err := os.Open("./Table/table.txt")
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

	asciiLength := 128
    stateLength := 21
    TT = make([][]int, stateLength)

    i := 0
    for scanner.Scan() {
        line := scanner.Text()
        values := strings.Fields(line)
        TT[i] = make([]int, asciiLength)

        for j, valueStr := range values {
            value, err := strconv.Atoi(valueStr)
            check(err)
            TT[i][j] = value
        }
        i++
        if i == stateLength {
            break
        }
    }

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading transition table:", err)
	}
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func New(input string) *Lexer {
    loadTransitionTable()
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() (token.Token, int) {
	var tok token.Token


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

			switch l.ch {
			case '%':
				tok = newToken(token.PERCENT)
				goto next
			case '*':
				tok = newToken(token.ASTERISK)
				goto next
			case '+':
				tok = newToken(token.PLUS)
				goto next
			case '-':
				tok = newToken(token.MINUS)
				goto next
			case '/':
				tok = newToken(token.SLASH)
				goto next
			case ';':
				tok = newToken(token.SEMICOLON)
				goto next
			case ',':
				tok = newToken(token.COMMA)
				goto next
			case '.':
				tok = newToken(token.EOF)
				goto next
			case '(':
				tok = newToken(token.LPAREN)
				goto next
			case ')':
				tok = newToken(token.RPAREN)
				goto next
			case '[':
				tok = newToken(token.LBRACE)
				goto next
			case ']':
				tok = newToken(token.RBRACE)
				goto next
			}
        case 1, 5, 8:
            l.readChar()
            l.skipWhitespace()
            state = TT[state][int(l.ch)]
		case 11, 13, 15, 16, 17, 18:
            str = str + string([]byte{l.ch})
            l.readChar()
            state = TT[state][int(l.ch)]
        case 19, 21:
            tok.Literal = str
            tok.Type = token.DOUBLE
            l.backChar()
            goto next
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
        case 6:
            tok = newToken(token.EQ)
            goto next
        case 7:
            tok = newToken(token.ASSIGN)
            l.backChar()
            goto next
        case 9:
            tok = newToken(token.GE)
			goto next
        case 10:
            tok = newToken(token.GT)
            l.backChar()
			goto next
        case 12:
            tok.Literal = str
            tok.Type = token.LookupIdent(tok.Literal)
            l.backChar()
            goto next
        case 14:
            str = str + string([]byte{l.ch})
            l.readChar()
            l.skipWhitespace()
            state = TT[state][int(l.ch)]
        case 20:
            tok.Literal = str
            tok.Type = token.INT
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

func (l *Lexer) backChar() {
	if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.position -= 1
        l.readPosition -= 1
        l.ch = l.input[l.position]
    }
}

func newToken(tokenType token.TokenType) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenType)}
}
