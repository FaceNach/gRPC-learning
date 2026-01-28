package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 3000
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello root route \n")
		w.Write([]byte("Hello root route"))
		fmt.Println("Hello root route")
	})
	
	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request){
		
		switch r.Method{
			case http.MethodGet:
				w.Write([]byte("Hello GET method on teachers route"))
				fmt.Println("Hello GET method on teachers route")
				return
			case http.MethodPost:
				w.Write([]byte("Hello POST method on teachers route"))
				fmt.Println("Hello POST method on teachers route")
				return
			case http.MethodPut:
				w.Write([]byte("Hello PUT method on teachers route"))
				fmt.Println("Hello PUT method on teachers route")
				return
			case http.MethodDelete:
				w.Write([]byte("Hello DELETE method on teachers route"))
				fmt.Println("Hello DELETE method on teachers route")
			return	
		}
	})
	
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello students route \n")
		w.Write([]byte("Hello students route"))
		fmt.Println("Hello students route")
	})
	
	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello execs route \n")
		w.Write([]byte("Hello execs route"))
		fmt.Println("Hello execs route")
	})
	
	fmt.Println("Listening server on port:", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil{
		log.Fatalln("Error creating the server: ", err)
	}
}

