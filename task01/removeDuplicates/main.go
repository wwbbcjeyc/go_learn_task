package main

import "fmt"

/*
*26. 删除有序数组中的重复项
给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，
返回删除后数组的新长度。元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。
考虑 nums 的唯一元素的数量为 k。去重后，返回唯一元素的数量 k。
nums 的前 k 个元素应包含 排序后 的唯一数字。下标 k - 1 之后的剩余元素可以忽略。
示例 1：
输入：nums = [1,1,2]
输出：2, nums = [1,2,_]
解释：函数应该返回新的长度 2 ，并且原数组 nums 的前两个元素被修改为 1, 2 。不需要考虑数组中超出新长度后面的元素。
示例 2：
输入：nums = [0,0,1,1,1,2,2,3,3,4]
输出：5, nums = [0,1,2,3,4,_,_,_,_,_]
解释：函数应该返回新的长度 5 ， 并且原数组 nums 的前五个元素被修改为 0, 1, 2, 3, 4 。不需要考虑数组中超出新长度后面的元素。
*/
func removeDuplicates(nums []int) int {

	result := make([]int, 0)

	result = append(result, nums[0])

	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			result = append(result, nums[i])
		}
	}
	return len(result)
}

// 双指针法
// 使用快慢指针，快指针遍历数组，慢指针指向当前不重复元素的位置。
func removeDuplicatesSZZ(nums []int) int {

	if (len(nums)) == 1 {
		return 1
	}

	slow := 1 // 慢指针，从第二个位置开始

	for fast := 1; fast < len(nums); fast++ {
		// 如果当前元素与前一个不同
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast] // 将不重复的元素移到前面
			slow++                  // 慢指针前进
		}
	}

	return slow

}

func main() {

	var nums []int = []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}

	fmt.Println(removeDuplicates(nums))

}
