/*
	[两数之和](https://leetcode-cn.com/problems/two-sum/)
	题目：给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
*/

package main

import "fmt"

func main() {
	variable1 := []int{2, 7, 11, 15}
	variable2 := []int{3, 2, 4}
	result1 := twoSum(variable1, 9)
	result2 := twoSum(variable2, 6)
	fmt.Println(result1)
	fmt.Println(result2)
}

func twoSum(nums []int, target int) []int {
	hashTable := make(map[int]int)
	for i, v := range nums {
		if key, is_exist := hashTable[target-v]; is_exist {
			return []int{key, i}
		}
		hashTable[v] = i

	}
	return nil
}
