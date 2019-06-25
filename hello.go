package main

import "fmt"

type Employee struct {
	firstName string
	lastName  string
	salary    int
	fullTime  bool
}

func main() {
	var roos Employee
	roos.firstName = "ross"
	roos.lastName = "Bing"
	roos.salary = 1200
	roos.fullTime = true

	fmt.Println("Hello World!, ", roos.firstName)
	fmt.Println("roos.firstName = ", roos.firstName)
	fmt.Println("roos.lastName = ", roos.lastName)
	fmt.Println("roos.salary = ", roos.salary)
	fmt.Println("roos.fullTime = ", roos.fullTime)
}
