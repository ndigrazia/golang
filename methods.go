package main

import (
	"fmt"
	"math"
)

/*Since there is no class-object architecture and the closest thing to class we have is structure.
Function with struct receiver is a way to achieve methods in Go.*/

//What is a method
//As you know that struct field can also be a function, the concept of method will be very easy for you to understand. A method is nothing but a function, but it belongs to a certain type. A method is defined with different syntax than normal function. It required an additional parameter known as a receiver which is a type to which the method belongs. A method can access properties of the receiver it belongs to.

type Employee struct {
	FirstName, LastName string
}

type Employee2 struct {
	name   string
	salary int
}

func (e Employee) fullName() string {
	return e.FirstName + " " + e.LastName
}

//Same method name
//One major difference between function and method is many methods can have the same name while no two functions with the same name can be defined in a package.
type Rect struct {
	width  float64
	height float64
}

type Circle struct {
	radius float64
}

type Contact struct {
	phone, address string
}
type Employee3 struct {
	name    string
	salary  int
	contact Contact
}

//Methods on nested struct
func (e *Employee3) changePhone3(newPhone string) {
	e.contact.phone = newPhone
}

//Value receivers
func (r Rect) Area() float64 {
	return r.width * r.height
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

//Pointer receivers
func (e *Employee2) changeName(newName string) {
	(*e).name = newName
}

func (e *Employee2) changeName2(newName string) {
	e.name = newName
}

//Methods on nested struct
func (c *Contact) changePhone4(newPhone string) {
	c.phone = newPhone
}

//Promoted field
type Employee4 struct {
	name   string
	salary int
	Contact
}

func (e *Employee4) changePhone5(newPhone string) {
	e.phone = newPhone
}

// A method can accept both pointer and value
type Employee6 struct {
	name   string
	salary int
}

func (e *Employee6) changeName8(newName string) {
	e.name = newName
}

func (e Employee6) showSalary() {
	e.salary = 1500
	fmt.Println("Salary of e =", e.salary)
}

func main() {
	e := Employee{
		FirstName: "Ross",
		LastName:  "Geller",
	}

	fmt.Println(e.fullName())

	//Pointer receivers
	e1 := Employee2{
		name:   "Ross Geller",
		salary: 1200,
	}
	// e before name change
	fmt.Println("e before name change =", e1)

	// create pointer to `e`
	ep := &e1

	// change name
	ep.changeName("Monica Geller")

	// e after name change
	fmt.Println("e after name change =", e1)

	e2 := Employee2{
		name:   "Ross Geller",
		salary: 1200,
	}

	fmt.Println("e before name change =", e2)
	// change name
	e2.changeName("Monica Geller")
	// e after name change
	fmt.Println("e after name change =", e2)

	//Methods on nested struct
	e3 := Employee3{
		name:    "Ross Geller",
		salary:  1200,
		contact: Contact{"011 8080 8080", "New Delhi, India"},
	}
	// e before phone change
	fmt.Println("e before phone change =", e3)
	// change phone
	e3.changePhone3("011 1010 1222")
	// e after phone change
	fmt.Println("e after phone change =", e3)

	e4 := Employee3{
		name:   "Ross Geller",
		salary: 1200,
		contact: Contact{
			phone:   "011 8080 8080",
			address: "New Delhi, India",
		},
	}
	// e before phone change
	fmt.Println("e before phone change =", e4)
	// change phone
	e4.contact.changePhone4("011 1010 1222")
	// e after phone change
	fmt.Println("e after phone change =", e4)

	//Promoted field
	e5 := Employee4{
		name:   "Ross Geller",
		salary: 1200,
		Contact: Contact{
			phone:   "011 8080 8080",
			address: "New Delhi, India",
		},
	}
	// e before phone change
	fmt.Println("e before phone change =", e5)
	// change phone
	e5.changePhone5("011 1010 1222")
	// e after phone change
	fmt.Println("e after phone change =", e5)

	//Promoted methods
	//Like promoted fields on a struct, methods implemented by inner struct is available on parent struct
	e6 := Employee4{
		name:   "Ross Geller",
		salary: 1200,
		Contact: Contact{
			phone:   "011 8080 8080",
			address: "New Delhi, India",
		},
	}
	// e before phone change
	fmt.Println("e before phone change =", e6)
	// change phone
	e6.changePhone4("011 1010 1222")
	// e after phone change
	fmt.Println("e after phone change =", e6)

	// A method can accept both pointer and value
	//When a function has a value argument, it will only accept the value of the parameter. If you passed a pointer to the function which expects a value, it will not work. This is also true when function accepts pointer, you simply can not pass a value to it.
	//When it comes to a method, that’s not the case. We can define a method with value or pointer receive and call it on pointer or value.
	e7 := Employee6{
		name:   "Ross Geller",
		salary: 1200,
	}
	// e before change
	fmt.Println("e before change =", e7)
	// calling `changeName` pointer method on value
	e7.changeName8("Monica Geller")
	// calling `showSalary` value method on pointer
	(&e7).showSalary()
	// e after change
	fmt.Println("e after change =", e7)
}
