package main

import (
	"fmt"
	"strconv"
)

// Mon Intermediate

type MonExpression interface{}

type MonLet struct {
	MonBindings []MonBinding
	Body        MonExpression
}

type MonBinding struct {
	Name  string
	Value MonExpression
}

type MonVar struct {
	Name string
}

type MonInt struct {
	Value int
}

type MonIf struct {
	Cond MonExpression
	Then MonExpression
	Else MonExpression
}

type MonBinary struct {
	Op    string
	Left  MonExpression
	Right MonExpression
}

type MonWhile struct {
	Cnd  MonExpression
	Body MonExpression
}

type MonSet struct {
	Var string
	Exp MonExpression
}

type MonBegin struct {
	Exps []MonExpression
}

// Select Instructor
type Instructions struct {
	Instructs [][]string
}

func PrintLetExpr(letExpr MonLet) string {
	bindings := letExpr.MonBindings
	variable := bindings[0].Name
	val := bindings[0].Value
	return "(let ((" + variable + " " +  PrintMon(val) + "))" + PrintMon(letExpr.Body) + ")"
}

func PrintMon(mon MonExpression) string  {
	switch e := mon.(type) {
	case MonInt:
		return strconv.Itoa(e.Value)
		
	case MonVar:
		return e.Name
	case MonLet:
		return PrintLetExpr(e)
	case MonIf:
		return "(if " + PrintMon(e.Cond) + " " + PrintMon(e.Then) + " " + PrintMon(e.Else) + ")"
	case MonBinary:
		return "(" + e.Op + " " + PrintMon(e.Left) + " " + PrintMon(e.Right) + ")"
	case MonWhile:
		return "(while " + PrintMon(e.Cnd) + " " + PrintMon(e.Body) + ")"
	case MonSet:
		return "(set " + PrintMon(e.Var) + " " + PrintMon(e.Exp) + ")"
	case MonBegin:
		result := "(begin"
		for i := range e.Exps {
			result += " " + PrintMon(i)
		}
		result += ")"
		return result
	default:
		return "Unknown expression"
	}
}

func ToAnf(expr Expression, counter int) MonExpression {
	switch e := expr.(type) {
	case IntLiteral:
		return MonInt{Value: e.Value}
	case LetExpr:
		monBindings := make([]MonBinding, len(e.Bindings))
		for i, bind := range e.Bindings {
			monBindings[i] = MonBinding{
				Name:  bind.Name,
				Value: ToAnf(bind.Value, counter ),
			}
		}
		monBody := ToAnf(e.Body, counter)
		return MonLet{MonBindings: monBindings, Body: monBody}
	case Var:
		return MonVar{Name: e.Name}
	case IfExpr:
		tmp := "temp_" + strconv.Itoa(counter)
		tmpVar := MonVar{Name: tmp}
		cnd := ToAnf(e.Cond, counter + 1)
		thn := ToAnf(e.Then, counter + 1)
		els := ToAnf(e.Else, counter + 1)
		ifExp := MonIf{Cond: tmpVar, Then: thn, Else: els}
		binding := MonBinding{Name: tmp, Value: cnd}
		monBindings := make([]MonBinding, 1)
		monBindings[0] = binding 
		return MonLet{MonBindings: monBindings, Body: ifExp}
	case BinaryOp:
		left := ToAnf(e.Left, counter)
		right := ToAnf(e.Right, counter)
		return MonBinary{Op: e.Operator, Left: left, Right: right}
	case WhileExpr:
		cnd := ToAnf(e.Cnd, counter)
		body := ToAnf(e.Body, counter)
		return MonWhile{Cnd: cnd, Body: body}
	case SetExpr:
		variable := e.Name
		exp := ToAnf(e.Value, counter)
		return MonSet{Var: variable, Exp: exp}
	case BeginExpr:
		exps := make([]MonExpression, len(e.Exprs))
		for i := range e.Exprs {
			exps[i] = ToAnf(e.Exprs[i], counter)
		}
		return MonBegin{Exps: exps}
		
	default:
		return nil
	}
}

func SelectInstructions(expr MonExpression) Instructions {
	switch e := expr.(type) {
	case MonInt:
		instructions := make([][]string, 0)
		strnum := strconv.Itoa(e.Value)
		movinstruction := []string{"movq", strnum, "%rdi"}
		callinstruction := []string{"callq", "print_int"}

		instructions = append(instructions, movinstruction, callinstruction)
		return Instructions{Instructs: instructions}

	case MonVar:
		instructions := make([][]string, 0)
		movinstruction := []string{"movq", e.Name, "%rdi"}
		callinstruction := []string{"callq", "print_int"}

		instructions = append(instructions, movinstruction, callinstruction)
		return Instructions{Instructs: instructions}

	case MonLet:
		instructions := make([][]string, 0)
		binding := e.MonBindings[0]

		switch val := binding.Value.(type) {
		case MonInt:
			strnum := strconv.Itoa(val.Value)
			movinstruction := []string{"movq", strnum, binding.Name}
			instructions = append(instructions, movinstruction)
			bodyinstructions := SelectInstructions(e.Body)
			instructions = append(instructions, bodyinstructions.Instructs...)
			return Instructions{Instructs: instructions}
		case MonBinary:
			op := val.Op
			switch op {
			case "<":
				leftExpr := val.Left
				rightExpr := val.Right

				var leftVarName, rightValue string
				switch left := leftExpr.(type) {
				case MonVar:  
					leftVarName = left.Name
				default:
					fmt.Println("Unsupported left operand type")
					return Instructions{Instructs: [][]string{}}
				}
				switch right := rightExpr.(type) {
				case MonInt:
					rightValue = strconv.Itoa(right.Value)
				default:
					fmt.Println("Unsupported right operand type")
					return Instructions{Instructs: [][]string{}}
				}
				instructions := [][]string{
					{"cmpq", "$" + rightValue, leftVarName},
					{"setl", "%al"},
					{"movzbq", "%al", "%rsi"},}
				bodyInstructions := SelectInstructions(e.Body)
				instructions = append(instructions, bodyInstructions.Instructs...)
				return Instructions{Instructs: instructions}
			default:
				fmt.Println("Unsupported binary operator")
				return Instructions{Instructs: [][]string{}}
			}
				
		default:
			fmt.Println("Unsupported MonExpression in Let")
			return Instructions{Instructs: [][]string{}}
		}

	case MonIf:
		
		instructions := [][]string{{"cmpq $1, %rsi"}, {"je block_16"}, {"jmp block_17"}, {"block_16"}}
		thenInstructions := SelectInstructions(e.Then)
		elseInstructions := SelectInstructions(e.Else)
		instructions = append(instructions, thenInstructions.Instructs...)
		instructions = append(instructions, []string{"jmp conclusion:"})
		instructions = append(instructions, []string{"block_17"})
		instructions = append(instructions, elseInstructions.Instructs...)
		instructions = append(instructions, []string{"jmp conclusion"})

		return Instructions{Instructs: instructions}
		

	case MonBinary:
		op := e.Op
		switch op {
		case "<":
			
			instructions := make([][]string, 0)
		
			rightExpr := e.Right
			leftExpr := e.Left
			switch valr := rightExpr.(type) {
			case MonInt:
				switch vall := leftExpr.(type) {
				case MonInt:
					n:= strconv.Itoa(valr.Value)
					n2 := strconv.Itoa(vall.Value)
					mv := []string{"movq", n, "temp_m0"}
					cmp := []string{"cmpq", "temp_m0", n2}
					instructions = append(instructions, mv, cmp)
					return Instructions{Instructs: instructions}
				case MonVar:
					strnum := strconv.Itoa(valr.Value)
					cmpin := []string{"cmpq", strnum, vall.Name}
					instructions = append(instructions, cmpin)
					return Instructions{Instructs: instructions}
				default:
					fmt.Println("Unsupported binary op")
					return Instructions{Instructs: [][]string{}}
					
				}
			default:
				fmt.Println("Unsupported binary operator")
				return Instructions{Instructs: [][]string{}}
			}

		default:
			fmt.Println("Unsupported binary operator")
			return Instructions{Instructs: [][]string{}}
		}
	case MonWhile:
		cnd := SelectInstructions(e.Cnd)
		cndins := cnd.Instructs
		body := SelectInstructions(e.Body)
		bodyins := body.Instructs

		instructions := make([][]string, 0)
		jllabel := [][]string{{"label", "loop"}}
		jlbody := append(jllabel, bodyins...)
		ins := append(instructions, jlbody...)
		jlin := [][]string{{"jl", "loop"}}
		inscmp := append(ins, cndins...)
		cmpjl := append(inscmp, jlin...)

		return Instructions{Instructs: cmpjl}
	case MonBegin:
		instructions := make([][]string, 0)
		for i := range e.Exps {
			instr := SelectInstructions(e.Exps[i])
			instructions = append(instructions, instr.Instructs...)
		}
		return Instructions{Instructs: instructions}

	case MonSet:
		switch exp := e.Exp.(type) {
		case MonInt:
			strnum := strconv.Itoa(exp.Value)
			ins := [][]string{{"movq", strnum, e.Var}}
			return Instructions{Instructs: ins}
		case MonBinary:
			er := exp.Right
			var strnum string
			switch exp2 := er.(type) {
			case MonInt:
				strnum = strconv.Itoa(exp2.Value)
				ins := [][]string{{"movq", strnum, "%rax"}, {"addq", "%rax", e.Var}}
				return Instructions{Instructs: ins}
			default:
				fmt.Println("Unsupported right-hand expression for binary operator")
			}
		default:
			fmt.Println("Unsupported expression in MonSet")
		}
		return Instructions{Instructs: [][]string{}}

	default:
		// If no case matches, return empty instructions
		fmt.Println("Unsupported expression type")
		return Instructions{Instructs: [][]string{}}
	}
}

func PrintSelect(ins Instructions) {
	fmt.Println(ins.Instructs)
}


func main() {
	input := "(let ((i 0)) (if (< i 3) 2 3))"
	//input := "(let ((i (let ((d 4)) (+ 3 3)))) (if (< 2 3) 2 3))"
	ast, _ := Parse(input)
	mon := ToAnf(ast, 0)
	mon_ := PrintMon(mon)
	fmt.Println(mon_)
	ss := SelectInstructions(mon)
	PrintSelect(ss)
	
}

