package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	FirstName string `json:"name"`
	Age       int    `json:"age"`
	Address Address `json:"address"`
}

type Address struct{
	City string `json:"city"`
	State string `json:"state"`
}

func jsonLesson() {
	person1 := person{
		FirstName: "John",
		Age:       30,
	}

	//marshalling
	jsonData, err:=json.Marshal(person1)
	if err != nil{
		fmt.Println("Json err: ", err)
	}

	fmt.Println(string(jsonData))

	person2 := person{
		FirstName: "Jane",
		Age: 35,
		Address: Address{
			City: "Parana",
			State: "Hurlingam",
		},
	}

	jsonData1 , err := json.Marshal(person2)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(string(jsonData1))
}