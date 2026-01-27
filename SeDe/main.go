package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

func main() {
	
	user := User{Name: "Alice", Email: "alice@example.com"}
	jsonData, err := json.Marshal(user)
	if err != nil{
		log.Fatal(err)
	}
		
	fmt.Println(string(jsonData))
	
	var user1 User
	err = json.Unmarshal(jsonData, &user1)
	
	if err != nil{
		log.Fatal(err)
	}
	
	fmt.Println("User created from json data", user1)
	
	jsonData1 := `{"name": "Jhon", "email":"jhon@example.com"}`
	
	reader := strings.NewReader(jsonData1)
	
	decoder := json.NewDecoder(reader)
	
	var user2 User 
	err = decoder.Decode(&user2)
	if err != nil{
		log.Fatal(err)
	}
	
	fmt.Println(user2.Name)
	fmt.Println(user2.Email)
	
	var buf bytes.Buffer
	enconder := json.NewEncoder(&buf)
	err = enconder.Encode(user)
	if err != nil{
		log.Fatal(err)
	}
	
	fmt.Println("Encoded: ", buf.String())
	
}