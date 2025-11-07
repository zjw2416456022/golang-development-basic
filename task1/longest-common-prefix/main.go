/*
	[最长公共前缀](https://leetcode-cn.com/problems/longest-common-prefix/)
	题目：查找字符串数组中的最长公共前缀
*/

package main

func main() {
	s1 := []string{"flower", "flow", "flght"}
	s2 := []string{"dog", "racecar", "car"}
	result1 := longestCommonPrefix(s1)
	result2 := longestCommonPrefix(s2)
	println(result1)
	println(result2)
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	for i := 0; i < len(strs[0]); i++ {
		for j := 1; j < len(strs); j++ {
			//最长不大于数组第一个元素的长度，匹配则继续判断，不匹配直接返回
			if !(i < len(strs[0]) && strs[j][i] == strs[0][i]) {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}
