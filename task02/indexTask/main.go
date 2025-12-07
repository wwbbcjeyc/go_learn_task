package main

import "fmt"

/*
*
编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，
在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
*/
func increase(nums *int) {
	*nums += 10
}

/**
实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
*/

func doubleSlice(nums *[]int) {

	for i := 0; i < len(*nums); i++ {
		(*nums)[i] *= 2
	}

}

func main() {

	a := 5
	fmt.Println("更改前:", a)

	increase(&a)
	fmt.Println("更改后:", a)

	nums := []int{1, 2, 3, 4, 5, 7}

	fmt.Println("Slice更改前:", nums)

	doubleSlice(&nums)

	fmt.Println("Slice更改后:", nums)

}
