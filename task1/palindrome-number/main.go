/*
	[回文数](https://leetcode-cn.com/problems/palindrome-number/)
	题目：给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
	回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
	例如，121 是回文，而 123 不是。
*/

package main

func main() {
	x1 := 121
	x2 := 123
	result1 := isPalindrome(x1)
	result2 := isPalindrome(x2)
	println(result1)
	println(result2)
}

func isPalindrome(x int) bool {
	// if x < 0 || (x > 0 && x%10 == 0) { //x > 0：需检查 “符号位 + 零标志位”（有符号数）或 “进位标志位”（无符号数）
	if x < 0 || (x != 0 && x%10 == 0) { //x != 0：直接检查 “零标志位（ZF）”
		return false
	}
	firstHalfNumber := 0
	for x > firstHalfNumber {
		firstHalfNumber = firstHalfNumber*10 + x%10
		x /= 10
	}
	//后面几位反转和前面几位相等，如果是奇数需要再左移一位
	return x == firstHalfNumber || x == firstHalfNumber/10
}
