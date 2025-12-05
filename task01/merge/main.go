package main

import (
	"fmt"
	"sort"
)

/**
56.合并区间
以数组 intervals 表示若干个区间的集合，
其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。
示例 1：

输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
输出：[[1,6],[8,10],[15,18]]
解释：区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].
示例 2：

输入：intervals = [[1,4],[4,5]]
输出：[[1,5]]
解释：区间 [1,4] 和 [4,5] 可被视为重叠区间。
示例 3：

输入：intervals = [[4,7],[1,4]]
输出：[[1,7]]
解释：区间 [1,4] 和 [4,7] 可被视为重叠区间。
*/

func merge(intervals [][]int) [][]int {

	if len(intervals) <= 1 {
		return intervals
	}

	// 1. 按起始位置排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 2. 初始化结果，放入第一个区间
	result := [][]int{intervals[0]}

	// 3. 遍历后续区间
	for i := 1; i < len(intervals); i++ {
		// 获取当前区间
		current := intervals[i]
		// 获取结果中的最后一个区间
		last := result[len(result)-1]

		// 判断是否可以合并
		if current[0] <= last[1] {
			// 可以合并，取结束位置的最大值
			if current[1] > last[1] {
				last[1] = current[1]
			}
		} else {
			// 不能合并，添加到结果
			result = append(result, current)
		}
	}

	return result

}

func main() {

	var arr1 [][]int = [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Println(merge(arr1))

}
