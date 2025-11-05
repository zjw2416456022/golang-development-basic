/*
	[有效的括号](https://leetcode-cn.com/problems/valid-parentheses/)
	题目：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
*/

package main

func main() {
	s1 := "()[]{}"
	s2 := "([)]"
	result1 := isValid(s1)
	result2 := isValid(s2)
	println(result1)
	println(result2)
}

func isValid(s string) bool {
	n := len(s)
	if n%2 == 1 { //奇数直接返回不匹配
		return false
	}
	pairs := map[byte]byte{ //定义有效括号
		')': '(',
		'}': '{',
		']': '[',
	}
	stack := []byte{} //定义栈
	for i := 0; i < n; i++ {
		if pairs[s[i]] == 0 { //没遇到右边括号就入栈
			stack = append(stack, s[i]) //入栈
		} else { //遇到右边括号
			if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] { //如果栈为空还遇到右边括号或者最顶上的栈不能和当前括号匹配则直接返回不匹配
				return false
			}
			stack = stack[:len(stack)-1] //出栈
		}
	}
	return len(stack) == 0 //栈为空则匹配
}
