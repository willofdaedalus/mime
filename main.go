package main

import "fmt"

func main() {
	// 	l := parser.NewParser(`entity student ->
	// 	id: number {increment}
	// 	dob: text
	// 	age: number
	// 	created_at: timestamp
	// 	gender: text <> ["male", "female"]
	//
	// # override the payload which is entity student by default
	// # with this custom one
	// alter ref student.payload ->
	// 	gender: text
	// 	age: number
	// 	dob: text
	//
	// alter ref student.response ->
	// 	id: number
	// 	dob: text
	// 	age: number
	// 	created_at: timestamp
	// 	gender: text
	//
	// # very basic routing; might consider advanced
	// routes ->
	// 	GET /employees/:id -> self.id
	// 	POST /employees -> payload
	// 	POST /employees -> self
	// `)
	//
	// 	l.NextToken()
	fmt.Println("hello mime")
}
