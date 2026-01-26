package main

import (
	"fmt"
	"log"
	"net/http"
)


func httpServer() {
	
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		fmt.Fprintln(res, "Hello server!")
	})
	
	const port string = ":3000"
	
	fmt.Println("Listening on port:", port)
	err := http.ListenAndServe(port,nil);
	if err != nil{
		log.Fatalln("error starting server:" , err)
	}
	
}