package main

import (
	"bufio"
	"fmt"
	"os"
)

var s string
var match []int
var memo map[[3]int]string

func solve(l, r int, close bool) string {
	key := [3]int{l, r, 0}
	if close {
		key[2] = 1
	}
	if v, ok := memo[key]; ok {
		return v
	}
	if l > r {
		memo[key] = ""
		return ""
	}
	if s[l] == '(' && match[l] != -1 && match[l] <= r {
		j := match[l]
		removeInner := solve(l+1, j-1, false)
		rest := solve(j+1, r, close)
		keepInner := solve(l+1, j-1, true)
		optionRemove := removeInner + rest
		optionKeep := "(" + keepInner + ")" + rest
		if close {
			if optionKeep+")" < optionRemove+")" {
				memo[key] = optionKeep
			} else {
				memo[key] = optionRemove
			}
		} else {
			if optionKeep < optionRemove {
				memo[key] = optionKeep
			} else {
				memo[key] = optionRemove
			}
		}
		return memo[key]
	}
	res := string(s[l]) + solve(l+1, r, close)
	memo[key] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &s)
	n := len(s)
	match = make([]int, n)
	for i := range match {
		match[i] = -1
	}
	stack := []int{}
	for i, ch := range s {
		if ch == '(' {
			stack = append(stack, i)
		} else {
			if len(stack) > 0 {
				j := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				match[i] = j
				match[j] = i
			}
		}
	}
	memo = make(map[[3]int]string)
	ans := solve(0, n-1, false)
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
