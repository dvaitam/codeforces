package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	s := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&s[i])
	}
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&c[i])
	}
	const inf int64 = 1 << 60
	ans := inf
	for j := 0; j < n; j++ {
		left := inf
		for i := 0; i < j; i++ {
			if s[i] < s[j] && c[i] < left {
				left = c[i]
			}
		}
		right := inf
		for k := j + 1; k < n; k++ {
			if s[j] < s[k] && c[k] < right {
				right = c[k]
			}
		}
		if left < inf && right < inf {
			cost := left + c[j] + right
			if cost < ans {
				ans = cost
			}
		}
	}
	if ans == inf {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
