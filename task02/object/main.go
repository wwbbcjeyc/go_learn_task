package main

import (
	"fmt"
	"math"
)

/*
*

定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，
实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法
*/
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 定义矩形结构体
type Rectangle struct {
	Width, Height float64
}

// 矩形实现Perimeter方法
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 矩形实现Area方法
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 定义圆形结构体
type Circle struct {
	Radius float64
}

// 圆形实现Perimeter方法
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 圆形实现Area方法
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

/**
使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
*/

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int `json:"employee_id"`
}

// 为Employee实现PrintInfo方法
func (e Employee) PrintInfo() {
	fmt.Println("=== 员工信息 ===")
	fmt.Printf("姓名: %s\n", e.Name)
	fmt.Printf("年龄: %d岁\n", e.Age)
	fmt.Printf("工号: %d\n", e.EmployeeID)
	fmt.Println("===============")
}

func main() {
	fmt.Println("=== 图形面积和周长计算 ===")

	rect := Rectangle{
		Width:  5.0,
		Height: 3.0,
	}

	circle := Circle{
		Radius: 2.5,
	}

	fmt.Println("直接调用:")
	fmt.Printf("矩形面积: %.2f\n", rect.Area())
	fmt.Printf("矩形周长: %.2f\n", rect.Perimeter())
	fmt.Printf("圆形面积: %.2f\n", circle.Area())
	fmt.Printf("圆形周长: %.2f\n", circle.Perimeter())

	fmt.Println("=== 结构体组 ===")

	emp := Employee{
		Person{
			Name: "张三",
			Age:  28,
		},
		10001,
	}

	emp.PrintInfo()
}
