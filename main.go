package main

import (
	"willofdaedalus/mime/internal/engine/lexer"
)

func main() {
	input := `entity user ->
	active int {default:"1"}
end`

	l := lexer.New(input)
	l.RenderTokens()
}
