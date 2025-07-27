package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := append([]int(nil), a...)
		sort.Ints(b)
		pos := make(map[int]int, n)
		for i, v := range b {
			pos[v] = i + 1 // 1-based positions
		}
		dp := make([]int, n+1)
		maxLen := 0
		for _, v := range a {
			idx := pos[v]
			dp[idx] = dp[idx-1] + 1
			if dp[idx] > maxLen {
				maxLen = dp[idx]
			}
		}
		fmt.Fprintln(out, n-maxLen)
	}
}
