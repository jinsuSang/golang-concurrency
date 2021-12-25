package main

import "fmt"

type person struct {
	first string
}

func (p person) speak() {
	fmt.Println("Hello ", p.first)
}

type human interface {
	speak()
}

func saySomething(h human) {
	h.speak()
}

func main() {
	p1 := person{first: "jinsu"}

	saySomething(p1)
}
