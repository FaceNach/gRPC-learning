package main

import "fmt"

func clousure() {

	// i := adder()

	// fmt.Printf("The current value of i: %v\n", i())
	// fmt.Printf("The current value of i: %v\n", i())
	// fmt.Printf("The current value of i: %v\n", i())
	// fmt.Printf("The current value of i: %v\n", i())

	subtracter := func() func(int) int {

		countdown := 99
		return func(x int) int {
			countdown -= x
			return countdown
		}
	}()

	fmt.Println(subtracter((1)))
	fmt.Println(subtracter((1)))
	fmt.Println(subtracter((1)))
	fmt.Println(subtracter((1)))
	fmt.Println(subtracter((1)))

}

func adder() func() int {

	i := 0
	fmt.Printf("previus value of i: %v \n", i)

	return func() int {
		i++
		fmt.Println("Added 1 to i")

		return i

	}
}
