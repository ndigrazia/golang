package main

import (
	"fmt"
	"math"
)

//Type aliasing
type float float64

//iota
const (
	a = iota + 1 // a == 1
	_            // (implicitly _ = iota + 1 )
	b            // b == 3 (implicitly b = iota + 1 )
	c            // c == 4 (implicitly c = iota + 1 )
)

func main() {
	a, b := 145.8, 543.8
	c := math.Min(a, b)
	fmt.Println("minimum value is ", c)

	//Type aliasing
	var f float = 52.2
	fmt.Printf("f has value %v and type %T", f, f)

	var other float64 = 10.0

	if float64(f) == other {
		fmt.Printf("equals.")
	} else {
		fmt.Printf("NOT equals.")
	}

}
