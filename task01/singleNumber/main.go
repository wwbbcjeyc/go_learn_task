package main

import (
	"fmt"
	"sort"
)

/*
*
136. 只出现一次的数字
给你一个 非空 整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
*/

/*
*
第一个想到的
*/
func singleNumber(nums []int) int {

	sort.Sort(sort.IntSlice(nums))

	if nums[0] != nums[1] {
		return nums[0]
	}

	for i := 1; i < len(nums)-1; i++ {
		if nums[i] != nums[i-1] && nums[i] != nums[i+1] {
			return nums[i]
		}
	}

	if nums[len(nums)-1] != nums[len(nums)-2] {
		return nums[len(nums)-1]
	}

	return -1

}

/*
*
第二个想到的
*/
func singleNUmberMap(nums []int) int {

	var numsMap = make(map[int]int)

	for i := 0; i < len(nums); i++ {
		//numsMap[nums[i]]
		numsMap[nums[i]] = numsMap[nums[i]] + 1
		//numsMap[nums[i]]++

	}

	for k, v := range numsMap {
		if v == 1 {
			return k
		}
	}

	return -1
}

/*
*
没想到的解法，异或运算原理：

	a ^ a = 0（相同数字异或为0）
	a ^ 0 = a（任何数与0异或等于自身）
	异或满足交换律和结合律
*/
func singleNumberYh(nums []int) int {

	result := 0
	for _, num := range nums {
		result ^= num

	}
	return result

}

func main() {

	//var nums [5]int = [5]int{4, 1, 2, 1, 2} //4
	var nums = []int{4, 5, 2, 4, 5, 2, 1, 2, 1, 3} //3
	//nums := []int{3, 4, 5, 4, 5, 3, 2, 1, 2} //1

	//fmt.Println(singleNumber(nums))
	fmt.Println(singleNUmberMap(nums))

}
