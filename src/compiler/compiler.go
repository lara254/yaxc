package compiler

import (
	"fmt"
	"strconv"
	"strings"
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
	return "(let ((" + variable + " " + PrintMon(val) + "))" + PrintMon(letExpr.Body) + ")"
}

func PrintMon(mon MonExpression) string {
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
				Value: ToAnf(bind.Value, counter),
			}
		}
		monBody := ToAnf(e.Body, counter)
		return MonLet{MonBindings: monBindings, Body: monBody}
	case Var:
		return MonVar{Name: e.Name}
	case IfExpr:
		tmp := "temp_" + strconv.Itoa(counter)
		tmpVar := MonVar{Name: tmp}
		cnd := ToAnf(e.Cond, counter+1)
		thn := ToAnf(e.Then, counter+1)
		els := ToAnf(e.Else, counter+1)
		ifExp := MonIf{Cond: tmpVar, Then: thn, Else: els}
		binding := MonBinding{Name: tmp, Value: cnd}
		monBindings := make([]MonBinding, 1)
		monBindings[0] = binding
		return MonLet{MonBindings: monBindings, Body: ifExp}
	case BinaryOp:
		if isAtomic(e.Left) && isAtomic(e.Right) {
			return MonBinary{Op: e.Operator, Left: ToAnf(e.Left, counter), Right: ToAnf(e.Right, counter)}
		} else if !isAtomic(e.Left) && isAtomic(e.Right) {
			tmp := "temp_" + strconv.Itoa(counter+1)
			tmpVar := MonVar{Name: tmp}
			addition := MonBinary{Op: e.Operator, Left: tmpVar, Right: ToAnf(e.Right, counter)}
			return makeLet(e.Left, addition, tmp, counter+2)
		} else if isAtomic(e.Left) && !isAtomic(e.Right) {
			tmp := "temp_" + strconv.Itoa(counter+1)
			tmpVar := MonVar{Name: tmp}
			addition := MonBinary{Op: e.Operator, Left: tmpVar, Right: ToAnf(e.Left, counter)}
			return makeLet(e.Right, addition, tmp, counter+2)
		} else {
			tmp := "temp_" + strconv.Itoa(counter+1)
			tmp2 := "temp_" + strconv.Itoa(counter+2)
			tmpVar := MonVar{Name: tmp}
			tmpVar2 := MonVar{Name: tmp2}
			addition := MonBinary{Op: e.Operator, Left: tmpVar, Right: tmpVar2}
			monLetExp := makeLet(e.Right, addition, tmp2, counter+3)
			return makeLet(e.Left, monLetExp, tmp, counter+4)
		}
	case WhileExpr:
		cnd := ToAnf(e.Cnd, counter)
		body := ToAnf(e.Body, counter)
		return MonWhile{Cnd: cnd, Body: body}
	case SetExpr:
		variable := e.Name
		exp := ToAnf(e.Value, counter)
		return MonSet{Var: variable, Exp: exp}
	case BeginExpr:
		letExp := BeginToLet(e)
		return ToAnf(letExp, counter)
	default:
		return nil
	}
}

func SelectInstructions(expr MonExpression, n int) Instructions {
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
			bodyinstructions := SelectInstructions(e.Body, n)
			instructions = append(instructions, bodyinstructions.Instructs...)
			return Instructions{Instructs: instructions}
		case MonBinary:
			op := val.Op
			switch op {
			case "<":
				leftExpr := val.Left
				rightExpr := val.Right

				var leftValue, rightValue string

				switch left := leftExpr.(type) {
				case MonVar:
					leftValue = left.Name
				case MonInt:
					leftValue = strconv.Itoa(left.Value)
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
					{"cmpq", "$" + rightValue, leftValue},
					{"setl", "%al"},
					{"movq", "%al", binding.Name},
				}
				bodyInstructions := SelectInstructions(e.Body, n)
				instructions = append(instructions, bodyInstructions.Instructs...)
				return Instructions{Instructs: instructions}
			default:
				fmt.Println("Unsupported binary operator")
				return Instructions{Instructs: [][]string{}}
			}

		case MonIf:
			switch cnd := val.Cond.(type) {
			case MonVar:
				mv := genMov(cnd, val.Then)
				mvElse := genMov(cnd, val.Else)
				cmpq := genCmpq(1, cnd.Name)
				block := makeBlock(n)
				block2 := makeBlock(n + 1)

				instructions := [][]string{cmpq, {"je", block}, {"jmp", block2}, {block}, mv, {block2}, mvElse}
				instructionsBody := SelectInstructions(e.Body, n)
				instructionsExp := append(instructions, instructionsBody.Instructs...)

				return Instructions{Instructs: instructionsExp}
			default:
				fmt.Println("unsupported IF condition, Must be atomic")
				return Instructions{Instructs: [][]string{}}
			}

		default:
			fmt.Println("Unsupported MonExpression in Let")
			return Instructions{Instructs: [][]string{}}
		}

	case MonIf:
		instructions := [][]string{{"cmpq $1, %rsi"}, {"je block_16"}, {"jmp block_17"}, {"block_16"}}
		thenInstructions := SelectInstructions(e.Then, n)
		elseInstructions := SelectInstructions(e.Else, n)
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
					n := strconv.Itoa(valr.Value)
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

		case "+":
			if isAnfVar(e.Left) && (isAnfVar(e.Right)) {
				return genAddition(e.Left, e.Right)
			} else {
				return Instructions{Instructs: [][]string{}}
			}

		default:
			fmt.Println("Unsupported binary operator")
			return Instructions{Instructs: [][]string{}}
		}

	case MonWhile:
		cnd := SelectInstructions(e.Cnd, n)
		cndins := cnd.Instructs
		body := SelectInstructions(e.Body, n+1)
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
			instr := SelectInstructions(e.Exps[i], n)
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
		fmt.Println("Unsupported expression type")
		return Instructions{Instructs: [][]string{}}
	}
}

func makeLet(expr MonExpression, expr2 MonExpression, tmp string, n int) MonExpression {
	anfExp := ToAnf(expr, n)
	switch ae := anfExp.(type) {
	case MonLet:
		bindings := ae.MonBindings
		letBody := ae.Body
		binding := MonBinding{Name: tmp, Value: letBody}
		bindings2 := make([]MonBinding, 1)
		bindings2[0] = binding
		return MonLet{MonBindings: bindings, Body: MonLet{MonBindings: bindings2, Body: expr2}}
	default:
		binding := MonBinding{Name: tmp, Value: expr}
		bindings := make([]MonBinding, 1)
		bindings[0] = binding
		return MonLet{MonBindings: bindings, Body: expr2}
	}
}

func BeginToLet(expr Expression) Expression {
	switch bgnExpr := expr.(type) {
	case BeginExpr:
		expsLength := len(bgnExpr.Exprs)
		last := bgnExpr.Exprs[expsLength-1]

		var result Expression = last

		for i := expsLength - 2; i >= 0; i-- {
			tmp := "tmp" + strconv.Itoa(i)
			binding := Binding{
				Name:  tmp,
				Value: bgnExpr.Exprs[i],
			}
			result = LetExpr{
				Bindings: []Binding{binding},
				Body:     result,
			}
		}
		return result
	default:
		return expr
	}
}

func isAtomic(expr Expression) bool {
	switch expr.(type) {
	case IntLiteral:
		return true
	case Var:
		return true
	default:
		return false
	}
}

func genAddition(e MonExpression, e2 MonExpression) Instructions {
	switch e1 := e.(type) {
	case MonVar:
		switch e3 := e.(type) {
		case MonVar:
			mov := genMov(e1, MonVar{Name: "%rax"})
			add := genAdd(MonVar{Name: "%rax"}, e3)

			return Instructions{Instructs: [][]string{mov, add}}
		default:
			return Instructions{Instructs: [][]string{}}
		}
	default:
		return Instructions{Instructs: [][]string{}}
	}
}
		
		
func genMov(cnd MonExpression, exp MonExpression) []string {
	switch cond := cnd.(type) {
	case MonVar:
		switch expr := exp.(type) {
		case MonInt:
			return []string{"movq", "$" + strconv.Itoa(expr.Value), cond.Name}
		case MonVar:
			return []string{"movq", cond.Name, expr.Name}
		default:
			return []string{}
		}
	default:
		return []string{"hello"}
	}
}

func genAdd(e MonExpression, e2 MonExpression) []string {
	switch exp := e.(type) {
	case MonVar:
		switch exp2 := e2.(type) {
		case MonVar:
			return []string{"addq", exp.Name, exp2.Name}
		default:
			return []string{}
		}
	default:
		return []string{}
	}
}

func isAnfVar(exp MonExpression) bool {
	switch exp.(type) {
	case MonVar:
		return true

	default:
		return false
	}
}
	
	


func genCmpq(boool int, cnd string) []string {
	return []string{"cmpq", "$" + strconv.Itoa(boool), cnd}
}

func makeBlock(block int) string {
	return "block" + strconv.Itoa(block)
}

func PrintSelect(ins Instructions) {
	fmt.Println(ins.Instructs)
}

func SelectInsToString(arr [][]string) string {
	rows := make([]string, len(arr))
	for i, row := range arr {
		rows[i] = "[" + strings.Join(row, " ") + "]"
	}
	return "[" + strings.Join(rows, " ") + "]"
}
