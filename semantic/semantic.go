package semantic

import (
    "fmt"
    "ezsharp/parser"
)

var debug = false

func typeError(line int, literal1 string, literal2 string) {
    fmt.Printf("Error:%d:Type mismatch '%s' and '%s'\n", line, literal1, literal2)
}

func peekSymbol(symbolTable [][]parser.Parse, currScope int) (parser.Parse) {
    if symbolTable[currScope] != nil {
        scopeLen := len(symbolTable[currScope])
        return symbolTable[currScope][scopeLen-1]
    } else {
        return parser.Parse{Token: "nil"}
    }
}

func Declare(symbolTable *[][]parser.Parse, curr parser.Parse, currScope int,
             tok string) {
    curr.Token = tok
    (*symbolTable)[currScope] = append((*symbolTable)[currScope], curr)
}

func Lookup(symbolTable [][]parser.Parse, token parser.Parse) (int, parser.Parse) {
	scopeNum := -1
    var returnNode parser.Parse

	for i := len(symbolTable) - 1; i >= 0; i-- {
		for j := len(symbolTable[i]) - 1; j >= 0; j-- {

            value := symbolTable[i][j]
            if value.Token == "FUNC_NAME" || value.Token == "VAR_RETURN" ||
               value.Token == "VAR_DECLARE" || value.Token == "VAR_ASSIGN" {
                if value.Literal == token.Literal {
                    scopeNum = i
                    returnNode = value
                    break
                }
			}
		}
	}
	return scopeNum, returnNode
}

func ReturnScope(symbolTable *[][]parser.Parse, curr *int) {
    (*symbolTable)[*curr] = nil
    *curr--
}

func NewScope(curr *int) {
	*curr++
}

func AnalyseSemantics(input []parser.Parse) ([]parser.Parse) {
    symbolTable := make([][]parser.Parse, 10)
	currScope := 0
    declType := ""
    space := 0
    flag := 0
    funcReturn := 0
    returnValue := ""

    output := make([]parser.Parse, 0)

if debug {
    println("Line:          Type:          Lexeme:             Token:          Scope:")
}

	for i := 0; i < len(input); i++ {
		curr := input[i]

        if curr.Type == "(" {
            if input[i+1].Token == "EXPR" {
            }

            if prev := peekSymbol(symbolTable, currScope); prev.Token != "nil" {
                if prev.Token == "FUNC_NAME" {
                    curr.Token = "SCOPE_START"
                } else if prev.Token == "FUNC_CALL" {
                    returnValue = prev.Type
                }
            }
        }

        if curr.Type == ")" {
            //NOTE: Possible error if bracket inside func paramaters
            if returnValue != "" {
                input[i].Type = returnValue
                returnValue = ""
            }
            funcReturn = 0
        }

		switch curr.Token {

        case "SCOPE_START":
            if curr.Literal == "else" {
                ReturnScope(&symbolTable, &currScope)
                space -= 3
            }
            NewScope(&currScope)
            space += 3
            break

        case "SCOPE_END":
            ReturnScope(&symbolTable, &currScope)
            space -= 3
            break

        case "INT_DECL":
            //declType = "INT_IDENT"
            declType = "INT"
            break

        case "DOUBLE_DECL":
            //declType = "DOUBLE_IDENT"
            declType = "DOUBLE"
            break

        case "FNAME":
            curr.Type = declType
            Declare(&symbolTable, curr, currScope, "FUNC_NAME")
            funcReturn = 1
            flag = 1
            break

        case "VAR":
            temp, vType := Lookup(symbolTable, curr)

            symbol := peekSymbol(symbolTable, currScope);

            if funcReturn == 1 {
                curr.Type = declType
                Declare(&symbolTable, curr, currScope, "VAR_RETURN")

            } else {
                if temp == -1 {
                    if symbol.Token == "VAR_DECLARE" || input[i-1].Literal == "int" ||
                       input[i-1].Literal == "double" {
                        curr.Type = declType
                        Declare(&symbolTable, curr, currScope, "VAR_DECLARE")
                        flag = 1
                        break
                    } else {
                        fmt.Printf("Error:%d:Unidentified identifier %s\n",
                            symbol.Line, curr.Literal)
                        goto error
                    }
                } else {
                    if input[i+2].Type == "INT" || input[i+2].Type == "DOUBLE" {
                        if vType.Type != input[i+2].Type {
                            typeError(symbol.Line+1, vType.Literal, input[i+2].Literal)
                            goto error
                        } else {
                            curr.Type = vType.Type
                            Declare(&symbolTable, curr, currScope, "VAR_ASSIGN")
                        }
                    } else if _, nType := Lookup(symbolTable, input[i+2]);
                        nType.Type != vType.Type {
                            typeError(symbol.Line+1, curr.Literal, nType.Literal)
                            goto error
                        } else {
                            curr.Type = vType.Type
                            Declare(&symbolTable, curr, currScope, "VAR_ASSIGN")
                        }
                    }
            }

            flag = 1
            break

        case "EXPR":
            if curr.Type == "INT" || curr.Type == "DOUBLE" {
                Declare(&symbolTable, curr, currScope, "LITERAL")
            }

            if curr.Type == "IDENT" {
                temp, definition := Lookup(symbolTable, curr)
                if temp == -1 {
                    if symbol := peekSymbol(symbolTable, currScope); symbol.Token != "nil" {
                        fmt.Printf("Error:%d:Unidentified identifier %s\n", symbol.Line, curr.Literal)
                        goto error
                    }
                }

                if definition.Token == "FUNC_NAME" {

                    if symbol := peekSymbol(symbolTable, currScope); symbol.Token != "nil" {
                        if symbol.Token == "VAR_ASSIGN" {
                            if symbol.Type == definition.Type {

                                curr.Type = definition.Type
                                Declare(&symbolTable, curr, currScope, "FUNC_CALL")
                            } else {
                                typeError(symbol.Line+1, symbol.Literal, definition.Literal)
                                goto error
                            }
                        } else if symbol.Token == "STATEMENT" || symbol.Token == "OPERATOR" {
                            curr.Type = definition.Type
                            Declare(&symbolTable, curr, currScope, "FUNC_CALL")
                        }
                    }
                }

                if definition.Token == "VAR_RETURN" ||
                   definition.Token ==  "VAR_ASSIGN" ||
                   definition.Token ==  "VAR_DECLARE" {
                    curr.Type = definition.Type
                    Declare(&symbolTable, curr, currScope, "VAR_REFERENCE")
                }
            }

            flag = 1
            break

        case "STATEMENT":
            Declare(&symbolTable, curr, currScope, curr.Token)
            flag = 1
            break

        case "RELOP":
            nType := input[i+1]
            pType := input[i-1]

            if input[i+1].Type == "IDENT" {
                _, nType = Lookup(symbolTable, input[i+1])
            }

            if input[i-1].Type == "IDENT" {
                _, pType = Lookup(symbolTable, input[i-1])
            }

            if nType.Type != pType.Type {
                typeError(curr.Line+1, pType.Literal, nType.Literal)
                goto error
            }

            Declare(&symbolTable, curr, currScope, curr.Token)
            flag = 1
            break

        case "OPERATOR":
            nType := input[i+1]
            pType := input[i-1]

            if nType.Literal == "(" {
                for j := i+1; j < len(input); j++ {
                    _, findFunc := Lookup(symbolTable, input[j])
                    if findFunc.Token == "FUNC_NAME" {
                        nType = input[j]
                        break
                    }
                }
            }

            if pType.Literal == ")" {
                for j := i-1; j > 0; j-- {
                    _, findFunc := Lookup(symbolTable, input[j])
                    if findFunc.Token == "FUNC_NAME" {
                        pType = input[j]
                        break
                    }
                }
            }

            if nType.Type == "IDENT" {
                _, nType = Lookup(symbolTable, nType)
            }

            if pType.Type == "IDENT" {
                _, pType = Lookup(symbolTable, pType)
            }

            if nType.Type != pType.Type {
                typeError(curr.Line+1, pType.Literal, nType.Literal)
                goto error
            }

            Declare(&symbolTable, curr, currScope, curr.Token)
            flag = 1
            break

        default:
            break
        }

        if curr.Literal == "then" {
            Declare(&symbolTable, curr, currScope, "then")
            flag = 1
        }

        if curr.Literal == "else" {
            Declare(&symbolTable, curr, currScope, "else")
            flag = 1
        }

        if curr.Literal == "fi" {
            Declare(&symbolTable, curr, currScope, "fi")
            flag = 1
        }

        if curr.Literal == "." {
            Declare(&symbolTable, curr, currScope, "EOF")
            flag = 1
        }

        if flag == 1 {
            if symbol := peekSymbol(symbolTable, currScope); symbol.Token != "nil" {
                symbol.Scope = currScope
                if currScope > 0 { fmt.Printf("%*s", symbol.Scope*3, "") }
                if debug {
                    fmt.Printf("{%-12d  %-16s  %-16s %-20s %10d}\n", symbol.Line, symbol.Type,
                    symbol.Literal, symbol.Token, symbol.Scope)
                }

                output = append(output, symbol)

                flag = 0
            }
        }
	}

    return output
    error:
    return nil
}
