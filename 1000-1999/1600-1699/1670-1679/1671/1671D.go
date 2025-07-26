package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solve(n, x int, a []int64) int64 {
	if n == 1 {
		v1 := int64(a[0] - 1)
		v2 := int64(x - 1)
		if v1 < v2 {
			return v2
		}
		return v1
	}
	base := int64(0)
	for i := 0; i < n-1; i++ {
		base += abs(a[i] - a[i+1])
	}
	costVal := func(val int64) int64 {
		res := min64(abs(a[0]-val), abs(a[n-1]-val))
		for i := 0; i < n-1; i++ {
			c := abs(a[i]-val) + abs(a[i+1]-val) - abs(a[i]-a[i+1])
			if c < res {
				res = c
			}
		}
		return res
	}
	ans := base + costVal(1)
	if x > 1 {
		ans += costVal(int64(x))
	}
	cand := base + int64(x-1) + min64(abs(a[0]-1), abs(a[0]-int64(x)))
	if cand < ans {
		ans = cand
	}
	cand = base + int64(x-1) + min64(abs(a[n-1]-1), abs(a[n-1]-int64(x)))
	if cand < ans {
		ans = cand
	}
	for i := 0; i < n-1; i++ {
		cand = base - abs(a[i]-a[i+1]) + int64(x-1) +
			min64(abs(a[i]-1)+abs(a[i+1]-int64(x)), abs(a[i]-int64(x))+abs(a[i+1]-1))
		if cand < ans {
			ans = cand
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, solve(n, x, a))
	}
}
