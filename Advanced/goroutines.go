package main

import (
	"fmt"
	"time"
)

func goroutines (){
	go sayHello()
}


func sayHello() {
	time.Sleep( 1 *time.Second)
	fmt.Println("Hello from sayHelloFunction")
}