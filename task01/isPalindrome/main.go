package main

import (
	"fmt"
	"strconv"
)

/*
*
9.回文数
给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。

回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。

例如，121 是回文，而 123 不是。

*/

/*
*
一开始想到：用字符串反转判断两个是否相等
*
*/
func isPalindrome(num int) bool {

	val := strconv.Itoa(num)

	runes := []rune(val)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	if val == string(runes) {
		return true
	}

	return false
}

/**
方法二：转换字符串
*/

func isPalindromeStr(num int) bool {

	str := strconv.Itoa(num)

	left := 0
	right := len(str) - 1
	for left < right {
		if str[left] != str[right] {
			return false
		}
		left++
		right--
	}

	return true
}

func main() {

	//s := 121
	//s := 123
	s := 454

	//fmt.Printf("是回文数嘛: %t\n", isPalindrome(s))

	fmt.Printf("是回文数嘛: %t\n", isPalindromeStr(s))

}
