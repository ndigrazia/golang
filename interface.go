package main

import (
	"fmt"
	"math"
	"strings"
)

//An interface is a collection of method signatures that an Object can implement
//But in Go, you do not explicitly mention if a type implements an interface

//Declaring interface

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rect struct {
	width  float64
	height float64
}

type Circle struct {
	radius float64
}

func (r Rect) Area() float64 {
	return r.width * r.height
}

func (r Rect) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

//Empty interface
func explain(i interface{}) {
	fmt.Printf("value given to explain function is of type '%T' with value %v\n", i, i)
}

type MyString string

//Multiple interfaces
type OtherShape interface {
	Area() float64
}

type Object interface {
	Volume() float64
}

type Cube struct {
	side float64
}

func (c Cube) Area() float64 {
	return 6 * (c.side * c.side)
}

func (c Cube) Volume() float64 {
	return c.side * c.side * c.side
}

type Shape3 interface {
	Area() float64
}

type Object3 interface {
	Volume() float64
}

type Skin3 interface {
	Color() float64
}

type Cube3 struct {
	side float64
}

func (c Cube3) Area() float64 {
	return 6 * (c.side * c.side)
}

func (c Cube3) Volume() float64 {
	return c.side * c.side * c.side
}

//Type switch
func explain1(i interface{}) {
	switch i.(type) {
	case string:
		fmt.Println("i stored string ", strings.ToUpper(i.(string)))
	case int:
		fmt.Println("i stored int", i)
	default:
		fmt.Println("i stored something else", i)
	}
}

//Embedding interfaces
//In Go, an interface cannot implement other interfaces or extend them, but we can create new interface by merging two or more interfaces.
type Material interface {
	OtherShape
	Object
}

func main() {
	//Implementing interface
	var s Shape
	s = Rect{5.0, 4.0}
	r := Rect{5.0, 4.0}
	fmt.Printf("type of s is %T\n", s)
	fmt.Printf("value of s is %v\n", s)
	fmt.Println("area of rectange s", s.Area())
	fmt.Println("s == r is", s == r)

	s = Circle{10}
	fmt.Printf("type of s is %T\n", s)
	fmt.Printf("value of s is %v\n", s)
	fmt.Printf("value of s is %0.2f\n", s.Area())

	// Empty interface
	ms := MyString("Hello World!")
	ms2 := "Hello World!"
	r1 := Rect{5.5, 4.5}

	explain(ms)
	explain(r1)
	explain(ms2)

	//Multiple interfaces
	c := Cube{3}
	var s1 OtherShape = c
	var o1 Object = c
	fmt.Println("volume of s of interface type Shape is", s1.Area())
	fmt.Println("area of o of interface type Object is", o1.Volume())

	//Type assertion
	var s2 Shape = Rect{3, 6}
	c2 := s2.(Rect)
	fmt.Println("area of c of type Cube is", c2.Area())
	fmt.Println("volume of c of type Cube is", c2.Perimeter())

	value, ok := s2.(Rect)
	fmt.Println("is ok", ok)
	fmt.Println("value", value)

	var s4 Shape3 = Cube3{3}
	value14, ok14 := s4.(Object3)
	fmt.Printf("dynamic value of Shape 's' with value %v implements interface Object? %v\n", value14, ok14)
	value24, ok24 := s4.(Skin3)
	fmt.Printf("dynamic value of Shape 's' with value %v implements interface Skin? %v\n", value24, ok24)

	//Type switch
	explain1("Hello World")
	explain1(52)
	explain1(true)

	//Embedding interfaces
	c6 := Cube{3}
	var m6 Material = c6
	var s6 OtherShape = c6
	var o6 Object = c6
	fmt.Printf("dynamic type and value of interface m of static type Material is'%T' and '%v'\n", m6, m6)
	fmt.Printf("dynamic type and value of interface s of static type Shape is'%T' and '%v'\n", s6, s6)
	fmt.Printf("dynamic type and value of interface o of static type Object is'%T' and '%v'\n", o6, o6)
}
