/*
	[136. 只出现一次的数字](https://leetcode.cn/problems/single-number/)
	题目：只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，
	其余每个元素均出现两次。找出那个只出现了一次的元素。可以使用 for 循环遍历数组，
	结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，
	然后再遍历 map 找到出现次数为1的元素。
*/

package main

func main() {
	nums := []int{4, 1, 2, 1, 2}
	result1 := singleNumber1(nums)
	result2 := singleNumber2(nums)
	println(result1)
	println(result2)
}
func singleNumber1(nums []int) int {
	result := 0
	countMap := make(map[int]int)
	//统计每个数字的出现次数
	for _, num := range nums {
		countMap[num]++
	}
	//找到出现次数为1的数字
	for num, count := range countMap {
		if count == 1 {
			result = num
		}
	}
	return result

}

func singleNumber2(nums []int) int {
	s := 0
	for _, num := range nums {
		s ^= num //转换成二进制按位异或,相同位为0，不同位为1
	}
	return s

}
