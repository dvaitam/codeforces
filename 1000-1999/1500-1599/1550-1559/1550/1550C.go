package main

import (
	"bufio"
	"fmt"
	"os"
)

func isGood(arr []int) bool {
	m := len(arr)
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			for k := j + 1; k < m; k++ {
				x := arr[i]
				y := arr[j]
				z := arr[k]
				if y >= min(x, z) && y <= max(x, z) {
					return false
				}
			}
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		var ans int64
		for l := 0; l < n; l++ {
			sub := make([]int, 0, 4)
			for r := l; r < n && r < l+4; r++ {
				sub = append(sub, a[r])
				if isGood(sub) {
					ans++
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
