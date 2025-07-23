package main

import (
	"bufio"
	"fmt"
	"os"
)

// lcp returns the length of the longest common prefix of a and b.
func lcp(a, b string) int {
	m := len(a)
	if len(b) < m {
		m = len(b)
	}
	for i := 0; i < m; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return m
}

// maxAreaWithLCP computes the maximum value of (r-l+1)*min(lcp[l..r-1])
// given the heights array which holds lcp of adjacent strings inside the range.
func maxAreaWithLCP(heights []int) int {
	n := len(heights)
	stack := make([]int, 0, n)
	maxArea := 0
	for i := 0; i <= n; i++ {
		var h int
		if i < n {
			h = heights[i]
		} else {
			h = -1 // sentinel to empty stack
		}
		for len(stack) > 0 && h < heights[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := -1
			if len(stack) > 0 {
				left = stack[len(stack)-1]
			}
			width := i - left
			area := width * heights[top]
			if area > maxArea {
				maxArea = area
			}
		}
		stack = append(stack, i)
	}
	return maxArea
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	s := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}
	for ; q > 0; q-- {
		var t int
		if _, err := fmt.Fscan(in, &t); err != nil {
			return
		}
		if t == 1 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			a--
			b--
			// max single string length
			ans := 0
			for i := a; i <= b; i++ {
				if len(s[i]) > ans {
					ans = len(s[i])
				}
			}
			if a < b {
				heights := make([]int, b-a)
				for i := a; i < b; i++ {
					heights[i-a] = lcp(s[i], s[i+1])
				}
				area := maxAreaWithLCP(heights)
				if area > ans {
					ans = area
				}
			}
			fmt.Fprintln(out, ans)
		} else if t == 2 {
			var x int
			var y string
			fmt.Fscan(in, &x, &y)
			s[x-1] = y
		}
	}
}
