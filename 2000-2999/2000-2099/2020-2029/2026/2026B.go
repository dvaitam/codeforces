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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		if n%2 == 0 {
			var ans int64
			for i := 0; i < n; i += 2 {
				diff := a[i+1] - a[i]
				if diff > ans {
					ans = diff
				}
			}
			fmt.Fprintln(out, ans)
			continue
		}

		left := make([]int64, n)
		for i := 1; i < n; i++ {
			if a[i]-a[i-1] == 1 {
				left[i] = left[i-1] + 1
			} else {
				left[i] = 0
			}
		}

		right := make([]int64, n)
		for i := n - 2; i >= 0; i-- {
			if a[i+1]-a[i] == 1 {
				right[i] = right[i+1] + 1
			} else {
				right[i] = 0
			}
		}

		const inf int64 = 1<<62 - 1
		ans := inf

		for skip := 0; skip < n; skip++ {
			var base int64
			var pending int64
			hasPending := false

			for j := 0; j < n; j++ {
				if j == skip {
					continue
				}
				if !hasPending {
					pending = a[j]
					hasPending = true
				} else {
					diff := a[j] - pending
					if diff > base {
						base = diff
					}
					hasPending = false
				}
			}

			g := left[skip] + 1
			if r := right[skip] + 1; r < g {
				g = r
			}
			cand := base
			if g > cand {
				cand = g
			}

			if cand < ans {
				ans = cand
			}
		}

		fmt.Fprintln(out, ans)
	}
}
