package main

import (
	"bufio"
	"fmt"
	"os"
)

type sub struct{ i, j int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(in, &n, &s)
		arr := []byte(s)
		// precompute LCP for all suffixes
		lcp := make([][]int, n+1)
		for i := 0; i <= n; i++ {
			lcp[i] = make([]int, n+1)
		}
		for i := n - 1; i >= 0; i-- {
			for j := n - 1; j >= 0; j-- {
				if arr[i] == arr[j] {
					lcp[i][j] = lcp[i+1][j+1] + 1
				}
			}
		}
		less := func(x, y sub) bool {
			len1 := x.j - x.i + 1
			len2 := y.j - y.i + 1
			l := lcp[x.i][y.i]
			if l >= len1 || l >= len2 {
				if len1 == len2 {
					return false
				}
				return len1 < len2
			}
			return arr[x.i+l] < arr[y.i+l]
		}
		tails := make([]sub, 0)
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				cur := sub{i, j}
				lo, hi := 0, len(tails)
				for lo < hi {
					mid := (lo + hi) / 2
					if less(tails[mid], cur) {
						lo = mid + 1
					} else {
						hi = mid
					}
				}
				if lo == len(tails) {
					tails = append(tails, cur)
				} else {
					tails[lo] = cur
				}
			}
		}
		fmt.Fprintln(out, len(tails))
	}
}
