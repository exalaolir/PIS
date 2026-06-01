package main

import "fmt"

type Celsius float64

func (c Celsius) String() string {
	return fmt.Sprintf("%.1f°C", c)
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) Birthday() {
	p.Age++
}

type Employee struct {
	Person
	Salary float64
}

type Greeter interface {
	Greet() string
}

func (p Person) Greet() string {
	return "Привет, я " + p.Name
}

func (e Employee) Greet() string {
	return "Сотрудник " + e.Name + ", зарплата " + fmt.Sprintf("%.2f", e.Salary)
}

func sayHello(g Greeter) {
	fmt.Println(g.Greet())
}

type People = []Person

func main() {

	var temp Celsius = 25.5
	fmt.Println(temp)

	p := Person{Name: "Андрей", Age: 30}
	p.Birthday()
	fmt.Println(p)

	emp := Employee{
		Person: Person{Name: "Вика", Age: 35},
		Salary: 50000,
	}
	fmt.Println(emp)
	fmt.Println("Возраст сотрудника:", emp.Age)
	emp.Birthday()
	fmt.Println("После дня рождения:", emp.Age)

	sayHello(p)
	sayHello(emp)

	var group People = []Person{{"Даша", 25}, {"Антон", 28}}
	fmt.Println(group)
}
