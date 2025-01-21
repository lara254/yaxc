package main

import (
	"fmt"
)

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
