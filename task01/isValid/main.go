package main

import "fmt"

/**20有效的括号
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。
输入：s = "()"  输出：true
输入：s = "()[]{}". 输出：true
输入：s = "(]".  输出：false
输入：s = "([])".  输出：true
输入：s = "([)]".  输出：false
*/

/**
使用栈（切片模拟）来存储左括号
遍历字符串中的每个字符
如果是左括号，就入栈
如果是右括号，检查栈顶是否匹配
最后检查栈是否为空
*/

func isValid(str string) bool {

	if len(str)%2 != 0 {
		return false
	}

	//切片构造栈
	stack := make([]byte, 0)

	//括号映射表
	paris := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	for i := 0; i < len(str); i++ {
		ch := str[i]

		// 右括号 进入匹配逻辑，判断成功，出栈
		if matchingLeft, isRight := paris[ch]; isRight {

			// 检查1：栈为空，没有左括号匹配
			// 检查2：栈顶元素不匹配
			if len(stack) == 0 || stack[len(stack)-1] != matchingLeft {
				return false
			}

			// 匹配成功，出栈
			stack = stack[:len(stack)-1]

		} else {
			// 左括号，入栈
			stack = append(stack, ch)
		}

	}
	return len(stack) == 0
}

func main() {
	//s := "()"
	//s := "()[]{}"
	//s := "(]"
	//s := "([])"
	s := "([)]"

	fmt.Printf("是否有效括号 %t\n", isValid(s))

	a := [5]int{6, 5, 4, 3, 2}
	fmt.Println(len(a))

	s7 := a[:len(a)-1]
	fmt.Println(s7)

}
