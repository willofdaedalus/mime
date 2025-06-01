package main

import "willofdaedalus/mime/internal/engine/lexer"

func main() {
	input := `entity user ->
	id uuid [primary required]
	name text [required unique]
	password text [hash]
	age int [default]
	end`

	l := lexer.New(input)
	l.RenderTokens()
}
