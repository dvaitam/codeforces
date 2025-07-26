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
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		// build initial breakpoints
		breaks := make([]int, 0, n)
		breaks = append(breaks, 1)
		for i := 1; i < n; i++ {
			if a[i] < a[breaks[len(breaks)-1]-1] {
				breaks = append(breaks, i+1)
			}
		}
		for j := 0; j < m; j++ {
			var k, d int
			fmt.Fscan(in, &k, &d)
			a[k-1] -= d
			pos := sort.SearchInts(breaks, k)
			idx := pos - 1
			prefix := int(1 << 60)
			if idx >= 0 {
				prefix = a[breaks[idx]-1]
			}
			has := pos < len(breaks) && breaks[pos] == k
			start := pos
			if has {
				if a[k-1] >= prefix {
					// remove breakpoint
					breaks = append(breaks[:pos], breaks[pos+1:]...)
				} else {
					start = pos + 1
				}
			} else {
				if a[k-1] < prefix {
					breaks = append(breaks, 0)
					copy(breaks[pos+1:], breaks[pos:])
					breaks[pos] = k
					start = pos + 1
				}
			}
			for start < len(breaks) && a[breaks[start]-1] >= a[k-1] {
				breaks = append(breaks[:start], breaks[start+1:]...)
			}
			if j+1 < m {
				fmt.Fprint(out, len(breaks), " ")
			} else {
				fmt.Fprintln(out, len(breaks))
			}
		}
	}
}
