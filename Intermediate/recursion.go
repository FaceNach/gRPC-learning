package main

import "fmt"

func recursion() {

	fmt.Println(factorial(5))
	fmt.Println(factorial(10))

	fmt.Println(sumOfDigits(9))
	fmt.Println(sumOfDigits(12))
	fmt.Println(sumOfDigits(12345))
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}

	return n * factorial(n-1)
}

func sumOfDigits(n int) int {
	if n < 10 {
		return n
	}

	return n%10 + sumOfDigits(n/10)
}
