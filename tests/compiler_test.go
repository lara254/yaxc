package main

import (
	"testing"
	"scmlike/src/compiler"
)

func TestAnf(t *testing.T) {
	ast, _ := compiler.Parse("(let ((i 0)) (if (< i 3) 2 3))")
	anf := compiler.ToAnf(ast, 0)
	got := compiler.PrintMon(anf)
	
	if got != "(let ((i 0))(let ((temp_0 (< i 3)))(if temp_0 2 3)))" {
		t.Errorf("ToAnf((let ((i 0)) (if (< i 3) 2 3))) = %s; want (let ((i 0))(let ((temp_0 (< i 3)))(if temp_0 2 3)))", got)
	}
}


func Testx86(t *testing.T) {
	ast, _ := compiler.Parse("(let ((i 0)) (if (< i 3) 2 3))")
	anf := compiler.ToAnf(ast, 0)
	ins := compiler.SelectInstructions(anf)
	got := ins.Instructs
	result := [][]string{
		{"movq", "0", "i"},
		{"cmpq", "$3", "i"},
		{"setl", "%al"},
		{"movzbq", "%al", "%rsi"},
		{"cmpq", "$1", "%rsi"},
		{"je", "block_16"},
		{"jmp", "block_17"},
		{"block_16"},
		{"movq", "2", "%rdi"},
		{"callq", "print_int"},
		{"jmp", "conclusion"},
		{"block_17"},
		{"movq", "3", "%rdi"},
		{"callq", "print_int"},
		{"jmp", "conclusion"},}

	for i := range len(got) {
		for j := range len(got[0]) {
			if got[i][j] != result[i][j] {
				t.Errorf("SelectInstructions((let ((i 0)) (if (< i 3) 2 3))) = %s; want %s", got[i][j], result[i][j])
			}
		}
	}
			
}
	
	
