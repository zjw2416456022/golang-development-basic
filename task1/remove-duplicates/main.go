/*
	[加一](https://leetcode-cn.com/problems/plus-one/)
	题目：给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
*/

package main

func main() {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 3, 2, 1}
	result1 := plusOne(s1)
	result2 := plusOne(s2)
	println(result1)
	println(result2)
}

func plusOne(digits []int) []int {
	lenght := len(digits)
	for i := lenght - 1; i >= 0; i-- {
		digits[i] += 1
		if digits[i] != 10 {
			return digits
		}
	}
	digits = make([]int, lenght+1)
	digits[0] = 1
	return digits
}
