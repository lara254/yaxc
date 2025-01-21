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
