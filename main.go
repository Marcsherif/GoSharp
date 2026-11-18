package main

import (
	"bufio"
	"errors"
	"ezsharp/lexer"
	"ezsharp/parser"
	"ezsharp/semantic"
	"ezsharp/codegen"
	"ezsharp/token"
	"fmt"
	"os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func createFile(path string) (*os.File) {
	file, err := os.Create(path)
	check(err)
	return file
}

func main() {
	lexerPass()
	parsedData, err := syntaxPass()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	symbolTable := semanticPass(parsedData)
	if symbolTable == nil {
		fmt.Fprintln(os.Stderr, "failed to get symbol table")
		return
	}
	IRPass(symbolTable)
}

func lexerPass() {
	// Pass an ezSharp source file to the compiler
	sourceFile, err := os.ReadFile(os.Args[1])
    check(err)

    lexerFile  := createFile("./output/lexerOutput.txt")
	errorFile  := createFile("./output/errors.txt")
    defer lexerFile.Close()
    defer errorFile.Close()

	tokens := bufio.NewWriter(lexerFile)
	errors := bufio.NewWriter(errorFile)

    l := lexer.New(string(sourceFile))
    for tok, line := l.NextToken(); true; tok, line = l.NextToken() {

        if tok.Type != token.ILLEGAL {
            t1 := fmt.Sprintf("{Type:%-8s Literal:%-8s line:%d}\n", tok.Type, tok.Literal, line + 1)
            tokens.WriteString(t1)
			parser.New(tok.Literal, line+1, string(tok.Type))
        } else {
            e1 := fmt.Sprintf("{Type:%-8s Literal:%-8s line:%d}\n", tok.Type, tok.Literal, line + 1)
            errors.WriteString(e1)
        }
        if tok.Type == token.EOF {
            break
        }
    }

	check(tokens.Flush())
	check(errors.Flush())
}

func syntaxPass() ([]parser.Parse, error) {
    parserFile := createFile("./output/parserOutput.txt")
    defer parserFile.Close()

    var parseOutput []parser.Parse
    var errOutput string

    parseOutput, errOutput = parser.SyntaxCheck()

    parserFile.WriteString("Type:          Lexeme:       Token:        Line:\n")
    for i := range parseOutput {
        parserFile.WriteString("-------------------------------------------------\n")
        t1 := fmt.Sprintf("%10s  %10s  %15s  %5d\n",
                   parseOutput[i].Type, parseOutput[i].Literal, parseOutput[i].Token, parseOutput[i].Line)
        parserFile.WriteString(t1)
    }

	if errOutput != "" {
		return nil, errors.New(errOutput)
	}

	return parseOutput, nil
}

func semanticPass(parserArr []parser.Parse) ([]parser.Parse) {
    symbolFile := createFile("./output/symbolOutput.txt")
    defer symbolFile.Close()

	symbolTable := semantic.AnalyseSemantics(parserArr)

    if symbolTable != nil {
        symbolFile.WriteString("Line:          Type:          Lexeme:             Token:          Scope:\n")
        for i := range symbolTable {
            if symbolTable[i].Scope > 0 {
				t1 := fmt.Sprintf("%*s", symbolTable[i].Scope*3, "");
				symbolFile.WriteString(t1)
			}

            t1 := fmt.Sprintf(
				"{%-12d  %-16s  %-16s %-20s %d}\n",
				symbolTable[i].Line, symbolTable[i].Type,
				symbolTable[i].Literal, symbolTable[i].Token,
				symbolTable[i].Scope,
			)
            symbolFile.WriteString(t1)
        }

		return symbolTable
    }

	return nil
}

func IRPass(symbolTable []parser.Parse) {
    tacFile    := createFile("./output/tacOutput.txt")
	defer tacFile.Close()

	t1 := codegen.IR(symbolTable)
	tacFile.WriteString(t1)
}
