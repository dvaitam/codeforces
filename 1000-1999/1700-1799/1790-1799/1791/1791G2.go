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
		var n int
		var c int64
		fmt.Fscan(in, &n, &c)
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		type pair struct {
			cost int64
			idx  int
		}
		pairs := make([]pair, n)
		for i := 1; i <= n; i++ {
			cost := int64(i)
			if tmp := int64(n + 1 - i); tmp < cost {
				cost = tmp
			}
			cost += a[i]
			pairs[i-1] = pair{cost, i}
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].cost < pairs[j].cost
		})
		pref := make([]int64, n)
		pos := make([]int, n+1)
		for i, p := range pairs {
			if i == 0 {
				pref[i] = p.cost
			} else {
				pref[i] = pref[i-1] + p.cost
			}
			pos[p.idx] = i
		}

		countWithBudget := func(b int64) int {
			if b < 0 {
				return 0
			}
			return sort.Search(n, func(i int) bool { return pref[i] > b })
		}

		ans := countWithBudget(c)
		for i := 1; i <= n; i++ {
			firstCost := int64(i) + a[i]
			if firstCost > c {
				continue
			}
			rem := c - firstCost
			k := countWithBudget(rem)
			if pos[i] < k {
				k--
			}
			if 1+k > ans {
				ans = 1 + k
			}
		}
		fmt.Fprintln(out, ans)
	}
}
