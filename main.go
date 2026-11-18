package main

import (
	"bufio"
	"errors"
	"ezsharp/codegen"
	"ezsharp/lexer"
	"ezsharp/parser"
	"ezsharp/semantic"
	"ezsharp/token"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	emitLLVM = true
	goSharpBuild = true
	goSharpRun = true
	sourcePath string
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
	flag.BoolVar(&emitLLVM, "ir", false, "keep the generated LLVM IR")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-ir] <source.cp>\n",
			filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	sourcePath = flag.Arg(0)

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
	tacPass(symbolTable)
	err = llvmPass(symbolTable)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func lexerPass() {
	sourceFile, err := os.ReadFile(sourcePath)
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

func tacPass(symbolTable []parser.Parse) {
    tacFile    := createFile("./output/tacOutput.txt")
	defer tacFile.Close()

	t1 := codegen.TacCodegen(symbolTable)
	tacFile.WriteString(t1)
}

func llvmPass(symbolTable []parser.Parse) (error) {
	build := goSharpBuild || goSharpRun
	if !emitLLVM && !build {
		return nil
	}

	irPath := "./output/llvmOutput.ll"
	exePath := "./output/objectOutput.exe"

	ir := codegen.LLVMCodegen(symbolTable)
	if err := os.WriteFile(irPath, []byte(ir), 0644); err != nil {
		return fmt.Errorf("write LLVM IR: %w", err)
	}

	if !build {
		return nil
	}

	cmd := exec.Command("clang", "-Wno-override-module", irPath, "-o", exePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clang failed: %w", err)
	}

	if goSharpRun {
		cmd := exec.Command(exePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("run compiled program: %w", err)
		}
	}

	if !emitLLVM {
		if err := os.Remove(irPath); err != nil {
			return fmt.Errorf("remove temporary LLVM IR: %w", err)
		}
	}


	return nil
}
