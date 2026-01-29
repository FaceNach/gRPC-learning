package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age int		`json:"age"`
	City string	`json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello root route \n")
		w.Write([]byte("Hello root route"))
		fmt.Println("Hello root route")
}

func teachersHandler(w http.ResponseWriter, r *http.Request){
	
	
		
		switch r.Method{
			case http.MethodGet:
				fmt.Println("URL params: ", r.URL.Path )
				path := strings.TrimPrefix(r.URL.Path, "/teachers/")
				userID := strings.TrimSuffix(path, "/")
				
				fmt.Println("URL id: ", userID)
				w.Write([]byte("Hello GET method on teachers route"))
				fmt.Println("Hello GET method on teachers route")
				fmt.Println("Query params: ", r.URL.Query())
				values := r.URL.Query()
				test := values.Get("test")
				phrase := values.Get("phrase")
				id := values.Get("id")
				
				fmt.Printf("Test: %v\nPhrase: %v\nId: %v\n", test, phrase, id)
				
			case http.MethodPost:
			
			//parse the form data(neccesary for form-urlencoded)
				err := r.ParseForm()
				if err != nil{
				http.Error(w, "Error parsing form", http.StatusBadRequest);
				return
			}
				fmt.Println("Form: ", r.Form)
				//prepare response data
				
				response := make(map[string]any)
				
				for key, value := range r.Form{
					response[key] = value[0]
				}	
				fmt.Println("Response map:" ,response)
				
				 body, err:= io.ReadAll(r.Body)
				if err != nil{
					return
				}
				
				defer r.Body.Close()
				
				fmt.Println(body)
				fmt.Println(string(body))
			
				var user User
				err = json.Unmarshal(body, &user )
				if err != nil {
					fmt.Println(err)
					return
				}
				
				fmt.Println("User: ", user)
				
				var user2  User
				err = json.NewDecoder(r.Body).Decode(&user2)
				
				fmt.Println("User2: ", user2)
				
				fmt.Println("Body: ", r.Body)
				fmt.Println("Form:" , r.Form)
				fmt.Println("Header:", r.Header)
				fmt.Println("Context: ", r.Context())
				fmt.Println("Content length: ", r.ContentLength)
				fmt.Println("Host: ", r.Host)
				fmt.Println("Protocol: ", r.Proto)
				fmt.Println("Remote addr: ", r.RemoteAddr)
				fmt.Println("Request URI: ", r.RequestURI)
				fmt.Println("TLS: ", r.TLS)
				fmt.Println("URL: ", r.URL)
				fmt.Println("User Agent: ", r.UserAgent())
				fmt.Println("Port: ", r.URL.Port())
				fmt.Println("Scheme: ", r.URL.Scheme)
				
				
				w.Write([]byte("Hello POST method on teachers route"))
				fmt.Println("Hello POST method on teachers route")
				
			case http.MethodPut:
				w.Write([]byte("Hello PUT method on teachers route"))
				fmt.Println("Hello PUT method on teachers route")
			
			case http.MethodDelete:
				w.Write([]byte("Hello DELETE method on teachers route"))
				fmt.Println("Hello DELETE method on teachers route")
		
		}
	}

func studentsHandler (w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello students route \n")
		w.Write([]byte("Hello students route"))
		fmt.Println("Hello students route")
		fmt.Println("Testing this air shit")
}

func execsHandler(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello execs route \n")
		w.Write([]byte("Hello execs route"))
		fmt.Println("Hello execs route")
	}


func main() {
	port := 3000
	
	http.HandleFunc("/", rootHandler)
	
	http.HandleFunc("/teachers/", teachersHandler)
	
	http.HandleFunc("/students", studentsHandler)
	
	http.HandleFunc("/execs", execsHandler)
	
	fmt.Println("Listening server on port:", port)
	
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil{
		log.Fatalln("Error creating the server: ", err)
	}
	

}

