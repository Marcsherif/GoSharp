package codegen

import (
	"strings"
	"ezsharp/parser"
	"fmt"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func LLVMCodegen(symbolTable []parser.Parse) (string) {
    var fileOutput strings.Builder
	fmt.Fprint(&fileOutput, "@fmt = private constant [4 x i8] c\"%d\\0A\\00\"\n\n")
	fmt.Fprint(&fileOutput, "declare i32 @printf(ptr, ...)\n\n")
	fmt.Fprint(&fileOutput, "define i32 @main() {\n")
	fmt.Fprint(&fileOutput, "entry:\n")

	var tempVarCount int

    for i := range symbolTable {
		symbol  := symbolTable[i]
		token   := symbol.Token
		literal := symbol.Literal

		switch (token) {
		case "VAR_DECLARE":
			fmt.Fprintf(&fileOutput, "  %%%s.addr = alloca i32\n", literal)
		case "VAR_ASSIGN":
			rhs, err := emitIntOperand(symbolTable[i+1], &fileOutput, &tempVarCount)
			check(err)
			fmt.Fprintf(&fileOutput, "  store i32 %s, ptr %%%s.addr\n", rhs, literal)
		case "STATEMENT":
			if literal != "print" {
				break
			}

			operand, err := emitIntOperand(symbolTable[i+1], &fileOutput, &tempVarCount)
			check(err)
			fmt.Fprintf(&fileOutput,
			"  call i32 (ptr, ...) @printf(ptr @fmt, i32 %s)\n",
			operand)
		}
	}

	fmt.Fprint(&fileOutput, "  ret i32 0\n}\n")
	return fileOutput.String()
}

func emitIntOperand(
	symbol parser.Parse, output *strings.Builder, tempCount *int,
) (string, error) {
	if symbol.Type != "INT" {
		return "", fmt.Errorf("expected an integer operand, got %q", symbol.Literal)
	}

	switch symbol.Token {
	case "LITERAL":
		tempName := fmt.Sprintf("%s", symbol.Literal)
		return tempName, nil
	case "VAR_REFERENCE":
		tempName := fmt.Sprintf("%%%s.value.%d", symbol.Literal, *tempCount)
		(*tempCount)++
		fmt.Fprintf(output, "  %s = load i32, ptr %%%s.addr\n", tempName, symbol.Literal)
		return tempName, nil
	default:
		return "", fmt.Errorf("unsupported operand %q (%s)", symbol.Literal, symbol.Token)
	}
}

