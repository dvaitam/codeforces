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

		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}

		pref := make([]int64, n)
		pref[0] = a[0]
		for i := 1; i < n; i++ {
			if a[i] > pref[i-1] {
				pref[i] = a[i]
			} else {
				pref[i] = pref[i-1]
			}
		}

		var ans int64
		for i := 0; i < n; i += 2 {
			var limit int64
			if i == 0 {
				limit = pref[1] - 1
			} else {
				// Odd positions (1-indexed) must stay strictly below the previous even peak.
				limit = pref[i-1] - 1
			}
			if a[i] > limit {
				ans += a[i] - limit
			}
		}
		fmt.Fprintln(out, ans)
	}
}
