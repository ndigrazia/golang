package main

import "fmt"

func getnumber() int {
	return 20
}

func main() {
	condition := true

	// if condition -- curly braces are mandatory
	if condition {
		fmt.Println("condition met")
	}

	//if - else condition
	a := 2

	if a > 10 {
		fmt.Println("condition met")
	} else {
		fmt.Println("condition did not meet")
	}

	// if — else if — else
	fruit := "orange"

	if fruit == "mango" {
		fmt.Println("fruit is mango")
	} else if fruit == "orange" {
		fmt.Println("fruit is orange")
	} else if fruit == "banana" {
		fmt.Println("fruit is banana")
	} else {
		fmt.Println("I don't know which fruit this is")
	}

	//Initial statement
	if ofruit := "banana"; ofruit == "mango" {
		fmt.Println("fruit is mango")
	} else if ofruit == "orange" {
		fmt.Println("fruit is orange")
	} else if ofruit == "banana" {
		fmt.Println("fruit is banana")
	} else {
		fmt.Println("I don't know which fruit this is")
	}

	// fruit variable is unavailable here
	//fmt.Println(ofruit)

	//Ternary condition
	//Go doesn’t support ternary one liners and there is no clear idiomatic way to achive that.

	//Switch: statement is used to replace multiple if-else conditions.
	//Unlike C, Go doesn't need a break statement to terminate the case block.
	//case block will terminate as soon as the last statement in the case block executes and so is the switch statement
	finger := 2

	switch finger {
	case 1:
		fmt.Println("Thumb")
	case 2:
		fmt.Println("Index")
	case 3:
		fmt.Println("Middle")
	case 4:
		fmt.Println("Ring")
	case 5:
		fmt.Println("Pinky")
	default:
		//default case does not have to be the last case. It can be anywhere in the switch block.
		fmt.Println("No fingers matched")
	}

	// multiple case values
	letter := "i"

	switch letter {
	case "a", "e", "i", "o", "u":
		fmt.Println("Letter is a vovel.")
	default:
		fmt.Println("Letter is not a vovel.")
	}

	//Initial statement statement
	//variable in the switch statement itself, restricting its scope to switch block.
	switch letter := "i"; letter {
	case "a", "e", "i", "o", "u":
		fmt.Println("Letter is a vovel.")
	default:
		fmt.Println("Letter is not a vovel.")
	}

	//Expressionless switch
	switch number := getnumber(); {
	case number <= 5:
		fmt.Println("number is less than or equal to 5")
	case number > 5:
		fmt.Println("number is greater than 5")
	case number > 10:
		fmt.Println("number is greater than 10")
	case number > 15:
		fmt.Println("number is greater than 15")
	}

	//fallthrough statement
	switch number := 20; {
	case number <= 5:
		fmt.Println("number is less than or equal to 5")
		fallthrough
	case number > 5:
		fmt.Println("number is greater than 5")
		fallthrough
	case number > 10:
		fmt.Println("number is greater than 10")
		fallthrough
	case number > 15:
		fmt.Println("number is greater than 15")
	}

	//For loops
	//Unlike C and other programming languages, the only available loop in Go is for loop. But many variants of for loop in Go will do all the jobs pretty well. So let's get into it.

	//for loop syntax
	//for loop syntax is consist of three statements and all three statements in Go are optional.
	for i := 1; i <= 6; i++ {
		fmt.Printf("Current number is %d \n", i)
	}

	//Optional post statement
	for i := 1; i <= 6; {
		fmt.Printf("Current number is %d \n", i)
		i++
	}

	//Optional init statement
	i := 1
	for ; i <= 6; i++ {
		fmt.Printf("Current number is %d \n", i)
	}

	//Optional init and post statement
	i := 1
	for i <= 6 { // ;i <= 6;
		fmt.Printf("Current number is %d \n", i)
		i++
	}

	//Without statements
	i := 1
	for {
		fmt.Printf("Current number is %d \n", i)

		if i == 6 {
			break
		}
		i++
	}

	//break statement
	//break statement is used from inside the for loop to terminate the for loop.
	for i := 1; i <= 10; i++ {
		if i > 6 {
			break
		}

		fmt.Printf("Current number is %d \n", i)
	}

	//continue statement
	//continue statement is used to skip one for loop iteration. When for loop sees continue statement, it simply ignores the current iteration, executes post statement and starts again
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			continue
		}

		fmt.Printf("Current number is %d \n", i)
	}

	fmt.Println("program terminated")
}
