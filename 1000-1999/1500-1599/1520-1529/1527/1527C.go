package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		// map from value to sum of indices of previous occurrences (1-indexed)
		sumIdx := make(map[int]int64)
		var ans int64
		for i := 1; i <= n; i++ {
			v := a[i-1]
			if prevSum, ok := sumIdx[v]; ok {
				ans += prevSum * int64(n-i+1)
			}
			sumIdx[v] += int64(i)
		}
		fmt.Fprintln(out, ans)
	}
}
