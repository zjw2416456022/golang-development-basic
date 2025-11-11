/*
	[两数56. 合并区间](https://leetcode.cn/problems/merge-intervals/)
	题目：以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
	请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
	可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，
	遍历排序后的区间数组，将当前区间与切片中最后一个区间进行比较，如果有重叠，则合并区间；
	如果没有重叠，则将当前区间添加到切片中。
*/

package main

import (
	"fmt"
	"sort"
)

func main() {
	variable1 := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	variable2 := [][]int{{1, 4}, {4, 5}}
	result1 := merge(variable1)
	result2 := merge(variable2)
	fmt.Println(result1)
	fmt.Println(result2)
}

func merge(intervals [][]int) [][]int {
	// if(len(intervals)==0){
	//     return [][]int{}
	// }
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	// merged:=[][]int{intervals[0]}
	merged := [][]int{}
	// for _,interval:= range intervals[1:] {
	for _, interval := range intervals {
		if len(merged) == 0 || merged[len(merged)-1][1] < interval[0] {
			merged = append(merged, interval)
		} else {
			if interval[1] > merged[len(merged)-1][1] {
				merged[len(merged)-1][1] = interval[1]
			}

		}
	}
	return merged
}
