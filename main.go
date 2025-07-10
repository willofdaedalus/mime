package main

import "willofdaedalus/mime/internal/dsl"

func main() {
	input := `
	app "todo"

	entity todo ->
	id: number [increment]
	item: text
	status: bool
	end

	@todo routes ->
	GET /todos/:id -> self.id
	POST /todos -> payload
	DELETE /todos/:id -> self.id
	`
	l := dsl.NewLexer(input)
	l.RenderTokens()
}
