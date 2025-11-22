package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type City struct {
	a, b int64
	d    int64
}

const inf int64 = 1<<62 - 1

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		cities := make([]City, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &cities[i].a, &cities[i].b)
			cities[i].d = cities[i].a - cities[i].b
		}
		sort.Slice(cities, func(i, j int) bool {
			if cities[i].d == cities[j].d {
				return cities[i].a < cities[j].a
			}
			return cities[i].d < cities[j].d
		})

		a := make([]int64, n)
		b := make([]int64, n)
		s := make([]int64, n)
		for i, c := range cities {
			a[i] = c.a
			b[i] = c.b
			s[i] = c.a + c.b
		}

		// valIncPrev: minimal (sum s - a_first) for subsequences of length len ending at i
		valIncPrev := make([]int64, n)
		valDecPrev := make([]int64, n) // minimal (sum s - b_first)
		for i := 0; i < n; i++ {
			valIncPrev[i] = b[i] // len=1: s_i - a_i = b_i
			valDecPrev[i] = a[i] // len=1: s_i - b_i = a_i
		}

		ans := make([]int64, n+1) // 1-based k
		for k := 2; k <= n; k++ {
			valIncCurr := make([]int64, n)
			valDecCurr := make([]int64, n)

			minInc := inf
			minDec := inf

			prefInc := valIncPrev[0]
			prefDec := valDecPrev[0]
			// i starts at 1 because need previous element
			for i := 1; i < n; i++ {
				valIncCurr[i] = s[i] + prefInc
				valDecCurr[i] = s[i] + prefDec

				if valIncCurr[i]-b[i] < minInc {
					minInc = valIncCurr[i] - b[i]
				}
				if valDecCurr[i]-a[i] < minDec {
					minDec = valDecCurr[i] - a[i]
				}

				if valIncPrev[i] < prefInc {
					prefInc = valIncPrev[i]
				}
				if valDecPrev[i] < prefDec {
					prefDec = valDecPrev[i]
				}
			}
			if minInc < minDec {
				ans[k] = minInc
			} else {
				ans[k] = minDec
			}

			valIncPrev = valIncCurr
			valDecPrev = valDecCurr
		}

		for k := 2; k <= n; k++ {
			if k > 2 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[k])
		}
		fmt.Fprintln(out)
	}
}
