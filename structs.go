package main

import (
	"fmt"
	"org"
)

//Creating a struct
type Employee struct {
	firstName string
	lastName  string
	salary    int
	fullTime  bool
}

// Anonymous fields
type Data struct {
	string
	int
	bool
}

//Nested struct
type Salary struct {
	basic     int
	insurance int
	allowance int
}

//Function fields
type FullNameType func(string, string) string

type Employee2 struct {
	firstName, lastName string
	salary              Salary
	bool
	FullName FullNameType
}

//struct field meta-data
type Employee6 struct {
	firstName string `json:"firstName"`
	lastName  string `json:"lastName"`
	salary    int    `json: "salary"`
	fullTime  int    `json: "fullTime"`
}

func main() {

	var roos Employee
	roos.firstName = "ross"
	roos.lastName = "Bing"
	roos.salary = 1200
	roos.fullTime = true

	//Getting and setting struct fields
	fmt.Println("roos.firstName = ", roos.firstName)
	fmt.Println("roos.lastName = ", roos.lastName)
	fmt.Println("roos.salary = ", roos.salary)
	fmt.Println("roos.fullTime = ", roos.fullTime)

	mary := Employee{
		firstName: "mary",
		lastName:  "Bing",
		salary:    1200,
		fullTime:  true,
	}

	fmt.Println(mary)

	//Initializing struct
	jhon := Employee{"jhon", "Bing", 1200, true}

	fmt.Println(jhon)

	//Anonymous struct
	monica := struct {
		firstName string
		lastName  string
		salary    int
		fullTime  bool
	}{
		firstName: "monica",
		lastName:  "Bing",
		salary:    1200,
		fullTime:  true,
	}

	fmt.Println(monica)

	//Pointer to struct
	peter := &Employee{
		firstName: "peter",
		lastName:  "Bing",
		salary:    1200,
		fullTime:  true,
	}

	fmt.Println("firstName", (*peter).firstName)
	fmt.Println("firstName", peter.firstName)

	// Anonymous fields
	sample1 := Data{"Monday", 1200, true}
	sample1.bool = false

	fmt.Println(sample1.string, sample1.int, sample1.bool)

	//Nested struct
	rossi := Employee2{
		firstName: "Rossi",
		lastName:  "Geller",
		bool:      true,
		salary:    Salary{1100, 50, 50},
	}
	fmt.Println(rossi)
	fmt.Println("Ross's basic salary", rossi.salary.basic)

	//Exported fields
	rovert := org.Employee3{
		FirstName: "Rovert",
		LastName:  "Geller",
	}
	fmt.Println(rovert)

	//Function fields
	e := Employee2{
		firstName: "Rossi",
		lastName:  "Geller",
		bool:      true,
		salary:    Salary{1100, 50, 50},
		FullName: func(firstName string, lastName string) string {
			return firstName + " " + lastName
		},
	}

	fmt.Println(e.FullName(e.firstName, e.lastName))

	//struct comparison
	e1 := Employee{
		firstName: "Ross",
		lastName:  "Geller",
		salary:    1200,
	}

	e2 := Employee{
		firstName: "Ross",
		lastName:  "Geller",
		salary:    1200,
	}

	fmt.Println(e1 == e2)

}
