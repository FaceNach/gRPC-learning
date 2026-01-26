package main

import (
	"fmt"
	"time"
)


func bufferedCh() {
	
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	
	go func () {
		time.Sleep(2 * time.Second)
		fmt.Println("Recieved: ", <-ch)
	}()
	fmt.Println("blocking ")
	ch <- 3
	fmt.Println("Recieved: ", <-ch)
	
	
	fmt.Println("End buffering")
	
}