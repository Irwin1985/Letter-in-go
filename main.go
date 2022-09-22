package main

import (
	parser2 "Letter/parser"
	"fmt"
)

func main() {
	input := `a.b.c.d.e.f.g[h].i.j.k.funcion()()()()()()(a, b, c)`
	parser := parser2.New(input)
	ast := parser.Parse()
	fmt.Println(ast.ToString())
}
