package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxQ = 300000 + 5

var (
	parent [maxQ]int
	amount [maxQ]int64
	cost   [maxQ]int64
	dsu    [maxQ]int
)

func find(x int) int {
	if dsu[x] != x {
		dsu[x] = find(dsu[x])
	}
	return dsu[x]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q, &amount[0], &cost[0])
	parent[0] = 0
	dsu[0] = 0

	for i := 1; i <= q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var p int
			var a, c int64
			fmt.Fscan(in, &p, &a, &c)
			parent[i] = p
			amount[i] = a
			cost[i] = c
			dsu[i] = i
		} else {
			var v int
			var w int64
			fmt.Fscan(in, &v, &w)
			bought := int64(0)
			spent := int64(0)
			for w > 0 {
				x := find(v)
				if amount[x] == 0 {
					if x == 0 {
						break
					}
					dsu[x] = find(parent[x])
					continue
				}
				take := w
				if take > amount[x] {
					take = amount[x]
				}
				amount[x] -= take
				w -= take
				bought += take
				spent += take * cost[x]
				if amount[x] == 0 && x != 0 {
					dsu[x] = find(parent[x])
				}
			}
			fmt.Fprintln(out, bought, spent)
		}
	}
}
