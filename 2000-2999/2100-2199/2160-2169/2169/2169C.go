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

		pref := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			pref[i] = pref[i-1] + v
		}

		const inf = int64(1 << 60)
		minG := inf
		best := int64(0)

		for i := 1; i <= n; i++ {
			ii := int64(i)
			// Derived gain uses G = (r^2 + r - pref[r]) - (l^2 - l - pref[l-1]).
			g := ii*ii - ii - pref[i-1]
			if g < minG {
				minG = g
			}
			candidate := ii*ii + ii - pref[i] - minG
			if candidate > best {
				best = candidate
			}
		}

		fmt.Fprintln(out, pref[n]+best)
	}
}
