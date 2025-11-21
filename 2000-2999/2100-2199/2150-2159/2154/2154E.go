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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &vals[i])
		}
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

		pref := make([]int64, n+1)
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i] + vals[i]
		}
		total := pref[n]
		bestDelta := int64(0)

		min := func(a, b int) int {
			if a < b {
				return a
			}
			return b
		}

		for j := 1; j <= n-2; j++ {
			maxM := min(j, n-1-j)
			if maxM <= 0 {
				continue
			}
			value := func(m int) int64 {
				v := vals[j]
				mk := int64(m) * int64(k)
				c := j
				if mk < int64(c) {
					c = int(mk)
				}
				loss := int64(m+1)*v - (pref[j+m+1] - pref[j])
				gain := int64(c)*v - pref[c]
				return loss + gain
			}

			low, high := 1, maxM
			for high-low > 5 {
				m1 := (2*low + high) / 3
				m2 := (low + 2*high) / 3
				v1 := value(m1)
				v2 := value(m2)
				if v1 < v2 {
					low = m1 + 1
				} else {
					high = m2 - 1
				}
			}
			for m := low; m <= high; m++ {
				val := value(m)
				if val > bestDelta {
					bestDelta = val
				}
			}
		}

		fmt.Fprintln(out, total+bestDelta)
	}
}
