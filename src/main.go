package main 

import (
	"fmt"
	"scmlike/src/compiler"
)

func main() {
	input := "(+ (if (< 2 3) 2 3) (if (< 3 4) 2 3))"
	//input := "(if (< 2 3) 3 4)"
	//input := "(< 2 3)"
	ast, _ := compiler.Parse(input)
	mon := compiler.ToAnf(ast, 0)
	//fmt.Println(mon)
	mon_ := compiler.PrintMon(mon)
	fmt.Println(mon_)
	
	/*
	ss := compiler.SelectInstructions(mon)
	compiler.PrintSelect(ss)
	*/
	
}
