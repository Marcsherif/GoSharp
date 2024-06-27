package main

import (
	"ezsharp/lexer"
	"ezsharp/token"
	"ezsharp/parser"
	"ezsharp/semantic"
	"ezsharp/tac"
	"fmt"
	"io"
	"os"
)

const bufferSize = 2048

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    f, err := os.Open(os.Args[1])
    check(err)

    defer f.Close()

    lexerf, err := os.Create("./output/lexerOutput.txt")
    check(err)

    parserf, err := os.Create("./output/parserOutput.txt")
    check(err)

    symbolf, err := os.Create("./output/symbolOutput.txt")
    check(err)

    tacf, err := os.Create("./output/tacOutput.txt")
    check(err)

    errorf, err := os.Create("./output/errors.txt")
    check(err)

    defer lexerf.Close()
    defer errorf.Close()

    buff1 := make([]byte, bufferSize)
    buff2 := make([]byte, bufferSize)
    var currentBuffer []byte
    var data []byte
    flag := true

    for flag == true {
        n, err := f.Read(buff1)
        if err == io.EOF {
            flag = false
        }
        if n > 0 {
            currentBuffer, buff1 = buff1, buff2
            buff2 = currentBuffer
            data = append(data, currentBuffer[:n]...)
        }
    }


    l := lexer.New(string(data))
    for tok, line := l.NextToken(); true; tok, line = l.NextToken() {

        if tok.Type != token.ILLEGAL {

            t1 := fmt.Sprintf("{Type:%-8s Literal:%-8s line:%d}\n", tok.Type, tok.Literal, line + 1)
            parser.New(tok.Literal, line+1, string(tok.Type))
            lexerf.WriteString(t1)

        } else {
            e1 := fmt.Sprintf("{Type:%-8s Literal:%-8s line:%d}\n", tok.Type, tok.Literal, line + 1)
            errorf.WriteString(e1)

        }
        if tok.Type == token.EOF {
            break
        }
    }

    var parserArr []parser.Parse
    var errOutput string

    parserArr, errOutput = parser.SyntaxCheck()

    parserf.WriteString("Type:          Lexeme:       Token:        Line:\n")
    for i := range parserArr {
        parserf.WriteString("-------------------------------------------------\n")
        t1 := fmt.Sprintf("%10s  %10s  %15s  %5d\n",
                   parserArr[i].Type, parserArr[i].Literal, parserArr[i].Token, parserArr[i].Line)
        parserf.WriteString(t1)
    }

    var symbolTable []parser.Parse
    if errOutput == "" {
        symbolTable = semantic.AnalyseSemantics(parserArr)
    } else {
        fmt.Printf("%+v\n", errOutput)
    }

    if symbolTable != nil {
        symbolf.WriteString("Line:          Type:          Lexeme:             Token:          Scope:\n")
        for i := range symbolTable {
            if symbolTable[i].Scope > 0 { t1 := fmt.Sprintf("%*s", symbolTable[i].Scope*3, ""); symbolf.WriteString(t1) }
            t1 := fmt.Sprintf("{%-12d  %-16s  %-16s %-20s %d}\n", symbolTable[i].Line, symbolTable[i].Type,
            symbolTable[i].Literal, symbolTable[i].Token, symbolTable[i].Scope)
            symbolf.WriteString(t1)
        }

        t1 := tac.IR(symbolTable)
        tacf.WriteString(t1)
    }

}
