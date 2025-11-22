package main

import (
	"bufio"
	"fmt"
	"os"
)

type query struct {
	t int
	i int
	x int
	k int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		qs := make([]query, q)
		gset := make(map[int]struct{})
		for idx := 0; idx < q; idx++ {
			var t int
			fmt.Fscan(in, &t)
			if t == 1 {
				var i, x int
				fmt.Fscan(in, &i, &x)
				qs[idx] = query{t: 1, i: i - 1, x: x}
			} else {
				var k int
				fmt.Fscan(in, &k)
				g := gcd(k, m)
				gset[g] = struct{}{}
				qs[idx] = query{t: 2, k: k}
			}
		}

		// Prepare list of needed gcds.
		glist := make([]int, 0, len(gset))
		for g := range gset {
			glist = append(glist, g)
		}

		// For each needed g, precompute number of descents in remainders.
		desc := make(map[int]int)
		for _, g := range glist {
			cnt := 0
			for i := 0; i+1 < n; i++ {
				if a[i]%g > a[i+1]%g {
					cnt++
				}
			}
			desc[g] = cnt
		}

		// Process queries in order.
		for _, qu := range qs {
			if qu.t == 1 {
				pos := qu.i
				oldVal := a[pos]
				newVal := qu.x
				for _, g := range glist {
					cnt := desc[g]
					if pos > 0 {
						left := a[pos-1] % g
						if left > oldVal%g {
							cnt--
						}
						if left > newVal%g {
							cnt++
						}
					}
					if pos+1 < n {
						right := a[pos+1] % g
						if oldVal%g > right {
							cnt--
						}
						if newVal%g > right {
							cnt++
						}
					}
					desc[g] = cnt
				}
				a[pos] = newVal
			} else {
				g := gcd(qu.k, m)
				if g == 1 {
					fmt.Fprintln(out, "YES")
					continue
				}
				cnt := desc[g]
				if cnt < m/g {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			}
		}
	}
}
