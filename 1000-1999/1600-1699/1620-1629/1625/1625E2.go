package main

import (
	"bufio"
	"fmt"
	"os"
)

func countRBS(arr []byte) int {
	stack := make([]int, 0, len(arr))
	dp := make([]int, len(arr))
	ans := 0
	for i, c := range arr {
		if c == '(' {
			stack = append(stack, i)
		} else if c == ')' {
			if len(stack) > 0 {
				open := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if open > 0 {
					dp[i] = dp[open-1] + 1
				} else {
					dp[i] = 1
				}
				ans += dp[i]
			}
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	sBytes := make([]byte, n)
	var s string
	fmt.Fscan(in, &s)
	copy(sBytes, []byte(s))

	for ; q > 0; q-- {
		var t, l, r int
		fmt.Fscan(in, &t, &l, &r)
		l--
		r--
		if t == 1 {
			if l >= 0 && l < n {
				sBytes[l] = '.'
			}
			if r >= 0 && r < n {
				sBytes[r] = '.'
			}
		} else {
			sub := sBytes[l : r+1]
			tmp := make([]byte, 0, len(sub))
			for _, ch := range sub {
				if ch == '(' || ch == ')' {
					tmp = append(tmp, ch)
				}
			}
			res := countRBS(tmp)
			fmt.Fprintln(out, res)
		}
	}
}
