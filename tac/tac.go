package tac

import (
	"ezsharp/parser"
	"fmt"
)

var debug = false
type OpCode string
var funcName []string

const (
    BG      OpCode = "BG"
    BL      OpCode = "BL"
    BE      OpCode = "BE"
    BGE     OpCode = "BGE"
    BLE     OpCode = "BLE"
    BNE     OpCode = "BNE"
    ADD     OpCode = "+"
    SUB     OpCode = "-"
    MUL     OpCode = "*"
    DIV     OpCode = "/"
    MOD     OpCode = "%"
)

type Quadruples struct {
    Op     OpCode
    Arg1   string
    Arg2   string
    Result string
}

func isFuncName(name string) bool {
    name += ":\n"
    for _, value := range funcName {
        if value == name {
            return true
        }
    }
    return false
}

func (q *Quadruples) String() string {
    return q.Result + " = " +  q.Arg1 + " " +  string(q.Op) + " " + q.Arg2
}

func GetRelop(relop string) OpCode {
    switch relop {
    case "<":
        return BL
    case "<=":
        return BLE
    case ">":
        return BG
    case ">=":
        return BGE
    case "==":
        return BE
    case "<>":
        return BNE
    default:
        return ""
    }
}

func GetOp(op string) OpCode {
    switch op {
    case "+":
        return ADD
    case "-":
        return SUB
    case "*":
        return MUL
    case "/":
        return DIV
    case "%":
        return MOD
    default:
        return ""
    }
}

func isOperator(op string) bool {
    switch op {
    case "+", "-", "*", "/", "%", "<", "<=", ">", ">=", "==", "<>":
        return true
    default:
        return false
    }
}

func priority(op string) int {
    switch op {
    case "+", "-":
        return 1
    case "*", "/", "%":
        return 2
    case "<", "<=", ">", ">=", "==", "<>":
        return 3
    default:
        return 0
    }
}

func infixToPostfix(formula []string) ([]string) {
    var stack []string
    var output []string

    for _, ch := range formula {
        if isOperator(ch) == false {
            output = append(output, ch)

        } else if ch == "(" {
            stack = append(stack, "(")

        } else if ch == ")" {

            for len(stack) > 0 && stack[len(stack)-1] != "(" {
                top := stack[len(stack)-1]
                stack = stack[:len(stack)-1]
                output = append(output, top)
            }
            stack = stack[:len(stack)-1]
        } else {

            for len(stack) > 0 && stack[len(stack)-1] != "(" && priority(ch) <= priority(stack[len(stack)-1]) {
                top := stack[len(stack)-1]
                stack = stack[:len(stack)-1]
                output = append(output, top)
            }
            stack = append(stack, ch)
        }
    }

    for len(stack) > 0 {
        top := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        output = append(output, top)
    }

    return output
}

func generateTAC(postfix []string, t *int) (string) {
    var exp_stack []string
    var tacReturn string

    //for x := len(postfix)-1; x >= 0; x-- {
    //    i := postfix[x]
    //    if isFuncName(i) == true {
    //        expStackLen := len(exp_stack)
    //        if len(exp_stack) > 0 {
    //            postfix[x+1] = exp_stack[expStackLen-1]
    //        }
    //        tacReturn += fmt.Sprintf("push {%s}\n", postfix[x+1])
    //        tacReturn += fmt.Sprintf("push {%s}\n", postfix[x+2])
    //        tacReturn += fmt.Sprintf("t%d = BL %s\n", *t, i)
    //        tacReturn += fmt.Sprintf("pop {%s}\n", postfix[x+2])
    //        tacReturn += fmt.Sprintf("pop {%s}\n", postfix[x+1])
    //        tacReturn += fmt.Sprintf("\n")
    //        if len(exp_stack) > 0 {
    //            exp_stack = exp_stack[:expStackLen-1]
    //        }
    //        postfix = append(postfix[:x], postfix[x+2:]...)
    //        exp_stack = append(exp_stack, fmt.Sprintf("t%d", *t))
    //        *t+=1
    //    }
    //}

    for x := 0; x < len(postfix); x++ {
        i := postfix[x]
        if isFuncName(i) == true {
            expStackLen := len(exp_stack)
            if len(exp_stack) > 0 {
                postfix[x+1] = exp_stack[expStackLen-1]
            }

            //var j, k int
            //var expr1, expr2 []string
            //for j = x+1; j < len(postfix); j++ {
            //    if postfix[j] == "," {
            //        expr1 = postfix[x+1:j]
            //        break
            //    }
            //}
            //for k = j+1; k < len(postfix); k++ {
            //    if isOperator(postfix[k]) || postfix[k] == "else" ||
            //        postfix[k] == "fi" {
            //        expr2 = postfix[j+1:k]
            //        break
            //    }
            //}
            expr1 := postfix[x+1]
            expr2 := postfix[x+2]

            tacReturn += fmt.Sprintf("push {%s}\n", expr1)
            tacReturn += fmt.Sprintf("push {%s}\n", expr2)
            tacReturn += fmt.Sprintf("t%d = BL %s\n", *t, i)
            tacReturn += fmt.Sprintf("pop {%s}\n", expr2)
            tacReturn += fmt.Sprintf("pop {%s}\n", expr1)
            tacReturn += fmt.Sprintf("\n")
            if len(exp_stack) > 0 {
                exp_stack = exp_stack[:expStackLen-1]
            }
            postfix = append(postfix[:x], postfix[x+3:]...)
            exp_stack = append(exp_stack, fmt.Sprintf("t%d", *t))
            *t+=1
        }
    }

    if postfix[0] == "return" {
        if len(postfix) == 2 {
            tacReturn += fmt.Sprintf("fp-4 = %s\n", exp_stack[0])
            return tacReturn
        }
    }
    if postfix[0] == "print" {
        if len(postfix) == 2 {
            tacReturn += fmt.Sprintf("print %s\n", exp_stack[0])
            return tacReturn
        }
    }

    for j, i := range postfix {
        if i == "return" {
            postfix[j] = "fp-4"
        }

        if len(postfix) == 3 {
            if len(exp_stack) > 0 {
                postfix[2] = exp_stack[0]
            }
			tacReturn += fmt.Sprintf("%s%s%s\n", postfix[0], postfix[1], postfix[2])
            break
        }

        if j == len(postfix)-1 {
            expStackLen := len(exp_stack)
			tacReturn += fmt.Sprintf("%s%s%s %s %s\n", postfix[0], postfix[1], exp_stack[expStackLen-2], i, exp_stack[expStackLen-1])
			exp_stack = exp_stack[:expStackLen-2]
            exp_stack = append(exp_stack, fmt.Sprintf("t%d", *t))
            break
        }

        if isOperator(i) == false {
            exp_stack = append(exp_stack, i)

        } else {
            expStackLen := len(exp_stack)
			tacReturn += fmt.Sprintf("t%d%s%s %s %s\n", *t, postfix[1], exp_stack[expStackLen-2], i, exp_stack[expStackLen-1])
			exp_stack = exp_stack[:expStackLen-2]
            exp_stack = append(exp_stack, fmt.Sprintf("t%d", *t))
			*t+=1
        }
    }

    tacReturn += fmt.Sprintf("\n")
    return tacReturn
}

func IR(symbolTable []parser.Parse) (string){
    funcCount := 0
    labCount := 0
    var exprCount int
    var postfix []string
    funcName = make([]string, 10)
    t := make([]int, 10)
    paramCount := make([]int, 10)
    varCount := make([]int, 10)
    tacCode := make([]string, 10)

    cmpCounter := 0
    cmpLocation := make([]int, 10)
    cmpFuncLocation := make([]int, 10)
    cmpArr := make([]string, 10)
    var fileOutput string

    if debug == true {
        fmt.Printf("B main\n")
    }
    fileOutput += fmt.Sprintf("B main\n")

    if symbolTable[0].Token != "FUNC_NAME" {
        if debug == true {
            fmt.Printf("main:\n")
            funcName[funcCount] += fmt.Sprintf("main:\n")
        }

        fileOutput += fmt.Sprintf("B main\n")
    }

    for i := range symbolTable {
        if symbolTable[i].Token == "FUNC_NAME" {
            if symbolTable[i].Literal != "main" {
                funcName[funcCount] += fmt.Sprintf("%s:\n", symbolTable[i].Literal)
                tacCode[funcCount] += fmt.Sprintf("push {LR}\npush {FP}\n")
            }
            if symbolTable[i].Literal == "main" {
                funcName[funcCount] += fmt.Sprintf("%s:\n", symbolTable[i].Literal)
            }
        }

        if symbolTable[i].Token == "VAR_RETURN" {
            paramCount[funcCount]++
        }

        if symbolTable[i].Token != "VAR_RETURN" {
            if i > 0 && symbolTable[i-1].Token == "VAR_RETURN" {
                for j := 1; j < paramCount[funcCount]+1; j++ {
//                    if symbolTable[i-j].Literal == "," { continue }
                    tacCode[funcCount] += fmt.Sprintf("%s = fp + %d\n", symbolTable[i-j].Literal, (j+1)*4)
                }
                tacCode[funcCount] += fmt.Sprintf("\n")
            }
        }

        if i > 0 && symbolTable[i].Scope == 0 {
            if symbolTable[i-1].Scope == 1 {
                tacCode[funcCount] += fmt.Sprintf("exit %s\n", funcName[funcCount])
                tacCode[funcCount] += fmt.Sprintf("pop {FP}\npop {PC}\n\n")

                paramCount[funcCount] = 0
                t[funcCount] = 0
                funcCount++

                if funcName[funcCount] != "main:\n" {
                    funcName[funcCount] += fmt.Sprintf("main:\n")
                }
            }
        }

        if symbolTable[i].Token == "VAR_ASSIGN" ||
           symbolTable[i].Literal == "print" ||
           symbolTable[i].Literal == "return" ||
           symbolTable[i].Literal == ";" {
            postfix = append(postfix, symbolTable[i].Literal)

            if symbolTable[i].Literal == "print" ||
               symbolTable[i].Token == "FUNC_CALL" {
                postfix = append(postfix, " ")
            } else {
                postfix = append(postfix, " = ")
            }

            for exprCount = i+1;
                symbolTable[exprCount].Token != "VAR_ASSIGN" &&
                symbolTable[exprCount].Token != "VAR_DECLARE" &&
                symbolTable[exprCount].Token != "STATEMENT" &&
                symbolTable[exprCount].Literal != "fi" &&
                symbolTable[exprCount].Literal != "else" &&
                symbolTable[exprCount].Literal != ".";
                exprCount++ {
                    postfix = append(postfix, symbolTable[exprCount].Literal)
            }

            pos := infixToPostfix(postfix)
            tacCode[funcCount] += generateTAC(pos, &t[funcCount])
            postfix = nil
        }

        if symbolTable[i].Token == "RELOP" {
            relop := GetRelop(symbolTable[i].Literal)
            cmpLocation[cmpCounter] = len(tacCode[funcCount])
            cmpArr[cmpCounter] += fmt.Sprintf("cmp %s, %s\n",
                symbolTable[i-1].Literal, symbolTable[i+1].Literal)
            cmpArr[cmpCounter] += fmt.Sprintf("%s lab%d\n", relop, labCount)
            cmpFuncLocation[cmpCounter] = funcCount
        }

        if symbolTable[i].Token == "then" {
            tacCode[funcCount] += fmt.Sprintf("lab%d:\n\n", labCount)
            labCount++
        }

        if symbolTable[i].Token == "else" {
            fName := funcName[funcCount][:len(funcName[funcCount])-2]
            tacCode[funcCount] += fmt.Sprintf("b exit%s\n", fName)
            tacCode[funcCount] += fmt.Sprintf("\n")
            cmpArr[cmpCounter] += fmt.Sprintf("b lab%d\n\n", labCount)
            tacCode[funcCount] += fmt.Sprintf("lab%d:\n\n", labCount)
            labCount++
        }

        if symbolTable[i].Token == "fi" {
            cmpCounter++
            fName := funcName[funcCount][:len(funcName[funcCount])-2]
            tacCode[funcCount] += fmt.Sprintf("b exit%s\n", fName)
            tacCode[funcCount] += fmt.Sprintf("\n")
        }

        if symbolTable[i].Token == "VAR_DECLARE" || symbolTable[i].Token == "VAR_RETURN" {
            varCount[funcCount]++
        }
    }

    var cmpLen int
    for f := 0; f < cmpCounter; f++ {
        funcLoc := cmpFuncLocation[f]
        if cmpLen > 0 {
            cmpLocation[f] += cmpLen
        }
        tacCode[funcLoc] = tacCode[funcLoc][:cmpLocation[f]] + cmpArr[f] + tacCode[funcLoc][cmpLocation[f]:]
        cmpLen = len(cmpArr[f])
    }

    for f := 0; f < funcCount+1; f++ {
//        tacCode[f] = tacCode[f][:cmpLocation[f]] + cmpArr[f] + tacCode[f][cmpLocation[f]:]
        if debug == true {
            fmt.Print(funcName[f])
            fmt.Printf("Begin: %d\n", (t[f]+varCount[f])*4)
            fmt.Print(tacCode[f])
        }
        fileOutput += fmt.Sprint(funcName[f])
        fileOutput += fmt.Sprintf("Begin: %d\n", (t[f]+varCount[f])*4)
        fileOutput += fmt.Sprint(tacCode[f])
    }
    return fileOutput
}

