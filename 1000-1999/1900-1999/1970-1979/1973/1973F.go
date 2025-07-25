package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ a, b int }

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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for tc := 0; tc < T; tc++ {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		b := make([]int, n)
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		ds := make([]int64, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &ds[i])
		}

		dp := map[pair]int64{}
		p0 := pair{a[0], b[0]}
		dp[p0] = 0
		alt := pair{b[0], a[0]}
		if v, ok := dp[alt]; !ok || int64(c[0]) < v {
			dp[alt] = int64(c[0])
		}

		for i := 1; i < n; i++ {
			ndp := make(map[pair]int64)
			for p, cost := range dp {
				g1 := gcd(p.a, a[i])
				g2 := gcd(p.b, b[i])
				key := pair{g1, g2}
				if v, ok := ndp[key]; !ok || cost < v {
					ndp[key] = cost
				}
				g1 = gcd(p.a, b[i])
				g2 = gcd(p.b, a[i])
				key = pair{g1, g2}
				cost2 := cost + int64(c[i])
				if v, ok := ndp[key]; !ok || cost2 < v {
					ndp[key] = cost2
				}
			}
			dp = ndp
		}

		type pv struct {
			cost int64
			val  int
		}
		arr := make([]pv, 0, len(dp))
		for p, cost := range dp {
			arr = append(arr, pv{cost, p.a + p.b})
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].cost < arr[j].cost })
		best := make([]int, len(arr))
		cur := 0
		for i, p := range arr {
			if p.val > cur {
				cur = p.val
			}
			best[i] = cur
		}

		for i, d := range ds {
			if i > 0 {
				out.WriteByte(' ')
			}
			j := sort.Search(len(arr), func(k int) bool { return arr[k].cost > d }) - 1
			if j >= 0 {
				fmt.Fprint(out, best[j])
			} else {
				fmt.Fprint(out, 0)
			}
		}
		out.WriteByte('\n')
	}
}
