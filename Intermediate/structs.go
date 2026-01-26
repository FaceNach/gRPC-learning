package main

import "fmt"



func structs() {


	p := Person {
		firstName: "Jhon",
		lastName: "Doe",
		age: 30,
	}

	p.fullName()
	printName(p.lastName)
	p.addOneYear()
}

type Person struct {
		firstName string
		lastName  string
		age       int
	}

func (p Person) fullName(){
	fmt.Println(p.firstName)
}

func (p *Person) addOneYear(){
	fmt.Println(p.age)
	p.age++
	fmt.Println(p.age)
}

func printName( n string) {
	fmt.Println(n)
}