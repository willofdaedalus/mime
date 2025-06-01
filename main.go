package main

import "willofdaedalus/mime/internal/engine/lexer"

func main() {
	input := `required]`

	l := lexer.New(input)
	l.RenderTokens()
}
