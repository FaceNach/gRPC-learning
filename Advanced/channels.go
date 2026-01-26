package main

import (
	"fmt"
	"time"
)

func channels() {
	greeting := make(chan string, 1)
	greetingString := "Hello"
	go func (phrase string){
		
		greeting <- phrase
		greeting <- "New Phrase"
		
		for _, e := range "abcde"{
			greeting <- "Alphabet: " + string(e)
		}
	}(greetingString)
	
	// go func(){
	// reciever := <- greeting
	// fmt.Println(reciever)
	
	// reciever = <- greeting
	// fmt.Println(reciever)
	// }()

	reciever := <- greeting
	fmt.Println(reciever)
	
	reciever = <- greeting
	fmt.Println(reciever)
	
	for range 5{
		rcv  := <- greeting
		fmt.Println(rcv)
	}
	
	time.Sleep(1 * time.Second)
	fmt.Println("End of program")
}