package parser

import (
	"fmt"
	"strings"
)

var debug = false

var nonTerminals []string
var terminals []string
var pArr []Parse

type Parse struct {
    Type         string
    Literal      string
    Token        string
    Line         int
    Scope        int
}

func New(literal string, line int, Type string) {
    pArr = append(pArr, Parse{Literal:literal, Line:line, Type:Type})
}

func setToken(input []Parse, index int, top string) (string) {

        if top == "<expr>"  || top == "<factor>" {
            input[index].Token = "EXPR"
        }

        if top == "(" {
            input[index].Token = ""
        }

        if top == "%" || top == "+"|| top == "-" ||
           top == "*" || top == "/" {
            input[index].Token = "OPERATOR"
        }

        if top == "<statement>" {
            input[index].Token = "STATEMENT"
        }

        if top == "while" || top == "then" || top == "else" {
            input[index].Token = "SCOPE_START"
        }

        if top == "od" || top == "fed" || top == "fi" {
            input[index].Token = "SCOPE_END"
        }

        if top == "int" {
            input[index].Token = "INT_DECL"
        }

        if top == "double" {
            input[index].Token = "DOUBLE_DECL"
        }

        if top == "<var>" {
            input[index].Token = "VAR"
        }

        if top == "=" {
            input[index].Token = "ASSIGN"
        }

        if top == "<>" || top == "<" || top == ">" || top == "<=" ||
           top == ">=" || top == "==" {
            input[index].Token = "RELOP"
        }

        if top == "<fname>" {
            input[index].Token = "FNAME"
        }


        return input[index].Token
}

func parse(input []Parse) ([]Parse, string) {
	terminals = append(terminals, "<program>", "<fdecls>", "<fdecls_prime>", "<fdec>", "<params>", "<params_prime>", "<fname>", "<declarations>", "<declarations_prime>", "<decl>", "<type>", "<varlist>", "<varlist_prime>", "<statement_seq>", "<statement_seq_prime>", "<statement>", "<statement_if>", "<expr>", "<expr_prime>", "<term>", "<term_prime>", "<factor>", "<factor_prime>", "<exprseq>", "<exprseq_prime>", "<bexpr>", "<bexpr_prime>", "<bterm>", "<bterm_prime>", "<bfactor>", "<bfactor_paren>", "<comp>", "<var>", "<var_prime>", "<letter>", "<digit>", "<id>", "<id_prime>", "<number>")
	nonTerminals = append(nonTerminals, "%", "=", "+", "-", "*", "/", "<=", "<", ">", ">=", "==", "<>", ",", ";", "(", ")", "[", "]", ".", "az", "and", "def", "do", "double", "else", "fed", "fi", "if", "int", "not", "od", "or", "print", "return", "then", "while", "#", "IDENT", "INT", "DOUBLE", "$")

    input = append(input, Parse{Literal:"$"})

	var stack []string
	start_symbol := "<program>"

	stack = append(stack, "$")
	stack = append(stack, start_symbol)

	index := 0
	flag := 0

	for len(stack) != 0 {
		top := stack[len(stack)-1]
		curr := input[index].Literal

        input[index].Token = setToken(input, index, top)

//        fmt.Printf("TOP: %21s, CURR_INPUT: %5s,  Token: %10s,  Line %d\n",
//                  top, curr, input[index].Token, input[index].Line)

        if input[index].Type == "IDENT"  {
            curr = "az"
        }
        if input[index].Type == "INT" {
            curr = "INT"
        }
        if input[index].Type == "DOUBLE" {
            curr = "DOUBLE"
        }

        if top == curr {
			stack = stack[:len(stack)-1]
			index++
		} else {
			row := getIndex(top, terminals)
			col := getIndex(curr, nonTerminals)
//            println("Row: ", row, "Col", col)
			value := parsingTable[row][col]

			if value == "-1" {
				flag = 1
				break
			}

			if value != "@" {
                strArr := strings.Split(value, " ")

                strArr = reverse(strArr)

                stack = stack[:len(stack)-1]

                for _, element := range strArr {
                    stack = append(stack, element)
                }
            } else {
				stack = stack[:len(stack)-1]
			}
		}
	}

    if debug {
        println("Type:          Lexeme:       Token:        Line:")
        for i := range input {
            println("-------------------------------------------------")
            fmt.Printf("%10s  %10s  %15s  %5d\n",
            input[i].Type, input[i].Literal, input[i].Token, input[i].Line)
        }
    }

    var errOutput string
	if flag == 0 {
	} else {
        errOutput = fmt.Sprintf("Error:line %d: '%s'\n", input[index].Line,
                   input[index-1].Literal)
	}

    return input, errOutput
}

func getIndex(target string, list []string) int {
	for i, str := range list {
		if str == target {
			return i
		}
	}
	return -1

}

func reverse(numbers []string) []string {
	newNumbers := make([]string, 0, len(numbers))
	for i := len(numbers)-1; i >= 0; i-- {
		newNumbers = append(newNumbers, numbers[i])
	}
	return newNumbers
}

func SyntaxCheck() ([]Parse, string){
    return parse(pArr)
}

var parsingTable = [][]string{
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<fdecls> <declarations> <statement_seq> .",	"<fdecls> <declarations> <statement_seq> .",	"-1",	"<fdecls> <declarations> <statement_seq> .",	"-1",	"<fdecls> <declarations> <statement_seq> .",	"-1",	"-1",	"-1",	"<fdecls> <declarations> <statement_seq> .",	"<fdecls> <declarations> <statement_seq> .",	"-1",	"-1",	"-1",	"<fdecls> <declarations> <statement_seq> .",	"<fdecls> <declarations> <statement_seq> .",	"-1",	"<fdecls> <declarations> <statement_seq> .",	"-1",	"-1",	"-1",	"-1",	"<fdecls> <declarations> <statement_seq> ."},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<fdecls_prime>",	"<fdecls_prime>",	"-1",	"<fdec> ; <fdecls_prime>",	"-1",	"<fdecls_prime>",	"-1",	"-1",	"-1",	"<fdecls_prime>",	"<fdecls_prime>",	"-1",	"-1",	"-1",	"<fdecls_prime>",	"<fdecls_prime>",	"-1",	"<fdecls_prime>",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"@",	"-1",	"<fdec> ; <fdecls_prime>",	"-1",	"@",	"-1",	"-1",	"-1",	"@",	"@",	"-1",	"-1",	"-1",	"@",	"@",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"def <type> <fname> ( <params> ) <declarations> <statement_seq> fed",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<type> <var> <params_prime>",	"-1",	"-1",	"-1",	"-1",	"<type> <var> <params_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	", <params>",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<id>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<declarations_prime>",	"<declarations_prime>",	"-1",	"-1",	"-1",	"<decl> ; <declarations_prime>",	"-1",	"<declarations_prime>",	"-1",	"<declarations_prime>",	"<decl> ; <declarations_prime>",	"-1",	"-1",	"-1",	"<declarations_prime>",	"<declarations_prime>",	"-1",	"<declarations_prime>",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"@",	"-1",	"-1",	"-1",	"<decl> ; <declarations_prime>",	"-1",	"@",	"-1",	"@",	"<decl> ; <declarations_prime>",	"-1",	"-1",	"-1",	"@",	"@",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<type> <varlist>",	"-1",	"-1",	"-1",	"-1",	"<type> <varlist>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"double",	"-1",	"-1",	"-1",	"-1",	"int",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<var> <varlist_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	", <varlist>",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<statement> <statement_seq_prime>",	"<statement> <statement_seq_prime>",	"-1",	"-1",	"-1",	"-1",	"<statement> <statement_seq_prime>",	"<statement> <statement_seq_prime>",	"<statement> <statement_seq_prime>",	"<statement> <statement_seq_prime>",	"-1",	"-1",	"<statement> <statement_seq_prime>",	"-1",	"<statement> <statement_seq_prime>",	"<statement> <statement_seq_prime>",	"-1",	"<statement> <statement_seq_prime>",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"; <statement_seq>",	"-1",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"@",	"@",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"@",	"<var> = <expr>",	"-1",	"-1",	"-1",	"-1",	"@",	"@",	"@",	"if <bexpr> then <statement_seq> <statement_if>",	"-1",	"-1",	"@",	"-1",	"print <expr>",	"return <expr>",	"-1",	"while <bexpr> do <statement_seq> od",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"else <statement_seq> fi",	"-1",	"fi",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<term> <expr_prime>",	"-1",	"-1",	"-1",	"-1",	"<term> <expr_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<term> <expr_prime>",	"<term> <expr_prime>",	"-1"},
    {"-1",	"@",	"+ <term> <expr_prime>",	"- <term> <expr_prime>",	"-1",	"-1",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"-1",	"@",	"-1",	"@",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"@",	"@",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<factor> <term_prime>",	"-1",	"-1",	"-1",	"-1",	"<factor> <term_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<factor> <term_prime>",	"<factor> <term_prime>",	"-1"},
    {"% <factor> <term_prime>",	"@",	"@",	"@",	"* <factor> <term_prime>",	"/ <factor> <term_prime>",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"-1",	"@",	"-1",	"@",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"@",	"@",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"( <expr> )",	"-1",	"-1",	"-1",	"-1",	"<id> <factor_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<number>",	"<number>",	"-1"},
    {"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"( <exprseq> )",	"<var_prime>",	"-1",	"<var_prime>",	"<var_prime>",	"-1",	"<var_prime>",	"-1",	"<var_prime>",	"-1",	"<var_prime>",	"<var_prime>",	"<var_prime>",	"-1",	"-1",	"-1",	"<var_prime>",	"<var_prime>",	"-1",	"-1",	"<var_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<expr> <exprseq_prime>",	"@",	"-1",	"-1",	"-1",	"<expr> <exprseq_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<expr> <exprseq_prime>",	"<expr> <exprseq_prime>",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	", <exprseq>",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<bterm> <bexpr_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<bterm> <bexpr_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"or <bterm> <bexpr_prime>",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<bfactor> <bterm_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<bfactor> <bterm_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"and <factor> <bterm_prime>",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"@",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"( <bfactor_paren>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"not <bfactor>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<expr> <comp> <expr> )",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<bexpr> )",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<expr> <comp> <expr> )",	"<expr> <comp> <expr> )",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<=",	"<",	">",	">=",	"==",	"<>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"<id> <var_prime>",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"-1",	"@",	"[ <expr> ]",	"@",	"@",	"-1",	"@",	"-1",	"@",	"-1",	"@",	"@",	"@",	"-1",	"-1",	"-1",	"@",	"@",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"az",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"#",	"#",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"az",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1"},
    {"@",	"@",	"@",	"@",	"@",	"@",	"-1",	"@",	"@",	"-1",	"-1",	"-1",	"@",	"@",	"@",	"@",	"@",	"@",	"@",	"az",	"@",	"-1",	"@",	"-1",	"@",	"@",	"@",	"-1",	"-1",	"-1",	"@",	"@",	"-1",	"-1",	"@",	"-1",	"-1",	"-1",	"#",	"#",	"-1"},
    {"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"-1",	"INT",	"DOUBLE",	"-1"}}
