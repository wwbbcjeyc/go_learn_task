package main

import "fmt"

/*
*1 两数之和
给定一个整数数组 nums 和一个整数目标值 target，
请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。

你可以按任意顺序返回答案。
示例 1：

输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。
示例 2：

输入：nums = [3,2,4], target = 6
输出：[1,2]
示例 3：

输入：nums = [3,3], target = 6
输出：[0,1]
*/
//暴力解法（双重循环）
func twoSum(nums []int, target int) []int {

	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

func toSumMap(nums []int, target int) []int {

	/*for i := 0; i < len(nums); i++ {
		com := target - nums[i]

		if value, exists := numMap[com]; exists {
			return []int{value, i}

		}
		numMap[nums[i]] = i
	}*/

	// 创建一个哈希表，key是数字，value是索引
	numMap := make(map[int]int)

	for i, num := range nums {
		// 计算需要的补数
		complement := target - num

		// 检查补数是否在哈希表中
		if idx, found := numMap[complement]; found {
			return []int{idx, i}
		}

		// 将当前数字和索引存入哈希表
		numMap[num] = i
	}

	return nil
}

func main() {

	nums := []int{3, 2, 4}
	target := 6

	fmt.Println(twoSum(nums, target))
	fmt.Println(toSumMap(nums, target))

}
