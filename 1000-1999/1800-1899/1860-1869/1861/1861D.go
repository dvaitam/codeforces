package main

import (
	"bufio"
	"fmt"
	"os"
)

func lisLength(arr []int) int {
	dp := make([]int, 0)
	for _, v := range arr {
		l, r := 0, len(dp)
		for l < r {
			m := (l + r) / 2
			if dp[m] < v {
				l = m + 1
			} else {
				r = m
			}
		}
		if l == len(dp) {
			dp = append(dp, v)
		} else {
			dp[l] = v
		}
	}
	return len(dp)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		l := lisLength(arr)
		fmt.Fprintln(out, n-l)
	}
}
