package main 

import (
	"fmt"
	"scmlike/src/compiler"
)

func main() {
	input := "(let ((i 0)) (if (< i 3) 2 3))"
	//input := "(let ((i (let ((d 4)) (+ 3 3)))) (if (< 2 3) 2 3))"
	ast, _ := compiler.Parse(input)
	mon := compiler.ToAnf(ast, 0)
	mon_ := compiler.PrintMon(mon)
	fmt.Println(mon_)
	ss := compiler.SelectInstructions(mon)
	compiler.PrintSelect(ss)
	
}
