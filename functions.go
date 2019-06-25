package main

import "fmt"

func dosomething() {
	fmt.Println("Hello World!")
}

func greet(user string) {
	fmt.Println("Hello " + user)
}
func add(a int, b int) {
	c := a + b
	fmt.Println(c)
}

// shorthand parameters notation
func substract(a, b int) {
	c := a - b
	fmt.Println(c)
}

//Return value
func add2(a, b int) int64 {
	return int64(a + b)
}

//Multiple return values
func addMult(a, b int) (int, int) {
	return a + b, a * b
}

//Named return values
func addMult1(a, b int) (add int, mul int) {
	add = a + b
	mul = a * b

	return // necessary
}

//defer keyword
func sayDone() {
	fmt.Println("I am done")
}

func endTime(timestamp string) {
	fmt.Println("Program ended at", timestamp)
}

// Function as type
type CalcFunc func(int, int) int

func calc(a int, b int, f CalcFunc) int {
	r := f(a, b)
	return r
}

func add1(a, b int) int {
	return a + b
}

func subtract1(a, b int) int {
	return a - b
}

//Function as value (anonymous function)
var addAsValue = func(a int, b int) int {
	return a + b
}

//Immediately-invoked function

func main() {
	dosomething()
	greet("John Doe")
	add(1, 3)
	substract(3, 1)
	fmt.Println(add2(1, 3))

	addRes, multRes := addMult(2, 5)
	fmt.Println(addRes, multRes)

	_, multRes1 := addMult(2, 5)
	fmt.Println(multRes1)

	//defer keyword
	fmt.Println("main started")

	defer sayDone()

	fmt.Println("main finished")

	time := "1 PM"

	defer endTime(time)

	time = "2 PM"

	fmt.Println("doing something")
	fmt.Println("main finished")
	fmt.Println("time is", time)

	fmt.Println("defer keyword Stack")

	//defer keyword
	fmt.Println("Call one")

	defer greet("Greet one")

	fmt.Println("Call two")

	defer greet("Greet two")

	fmt.Println("Call three")

	defer greet("Greet three")

	addResult1 := calc(5, 3, add1)
	subResult1 := calc(5, 3, subtract1)
	fmt.Println("5+3 =", addResult1)
	fmt.Println("5-3 =", subResult1)

	fmt.Println("5+3 =", addAsValue(5, 3))

	//Immediately-invoked function
	sum := func(a int, b int) int {
		return a + b
	}(8, 9)

	fmt.Println("8+9 =", sum)
}
