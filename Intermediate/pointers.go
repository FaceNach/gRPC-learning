package main

import "fmt"

func pointers() {
	var ptr *int
	var a int = 10

	var nilPtr *string

	ptr = &a

	fmt.Println(a)
	fmt.Println(&a)
	fmt.Println(ptr) //direccion en memoria a lo que apunta
	fmt.Println(&ptr) //direccion en memoria del puntero
	fmt.Println(*ptr) //valor de lo que apunta en memoria

	if nilPtr == nil{
		fmt.Println(nilPtr)
	}
	

}