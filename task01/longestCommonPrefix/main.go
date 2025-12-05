package main

import "fmt"

/**
编写一个函数来查找字符串数组中的最长公共前缀。

如果不存在公共前缀，返回空字符串 ""。
示例 1：

输入：strs = ["flower","flow","flight"]
输出："fl"
示例 2：

输入：strs = ["dog","racecar","car"]
输出：""
解释：输入不存在公共前缀。
*/

// 横向法
func longestCommonPrefix(strs []string) string {

	if (len(strs)) == 0 {
		return ""
	}

	// 以第一个字符串为基准
	prefix := strs[0]

	for i := 1; i < len(strs); i++ {
		// 逐个字符比较，直到找到不匹配的位置
		j := 0
		for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
			j++
		}

		// 更新前缀
		prefix = prefix[:j]

		// 如果前缀已经为空，提前返回
		if prefix == "" {
			return ""
		}
	}

	return prefix
}
func main() {

	str := []string{"flower", "flow", "flight"}

	fmt.Println(longestCommonPrefix(str))
}
