package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int64
	idx int
}

const INF int64 = 1 << 60

func compute(costs []int64, next []int64, bans []map[int]struct{}) []int64 {
	m := len(costs)
	res := make([]int64, m)
	pairs := make([]pair, len(next))
	for i, v := range next {
		pairs[i] = pair{v, i}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
	for i := 0; i < m; i++ {
		best := INF
		b := bans[i]
		for _, p := range pairs {
			if p.val >= INF {
				break
			}
			if b != nil {
				if _, ok := b[p.idx]; ok {
					continue
				}
			}
			best = costs[i] + p.val
			break
		}
		res[i] = best
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n1, n2, n3, n4 int
	fmt.Fscan(in, &n1, &n2, &n3, &n4)
	a := make([]int64, n1)
	b := make([]int64, n2)
	c := make([]int64, n3)
	d := make([]int64, n4)
	for i := 0; i < n1; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n2; i++ {
		fmt.Fscan(in, &b[i])
	}
	for i := 0; i < n3; i++ {
		fmt.Fscan(in, &c[i])
	}
	for i := 0; i < n4; i++ {
		fmt.Fscan(in, &d[i])
	}

	ban12 := make([]map[int]struct{}, n1)
	var m1 int
	fmt.Fscan(in, &m1)
	for i := 0; i < m1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		if ban12[x] == nil {
			ban12[x] = make(map[int]struct{})
		}
		ban12[x][y] = struct{}{}
	}

	ban23 := make([]map[int]struct{}, n2)
	var m2 int
	fmt.Fscan(in, &m2)
	for i := 0; i < m2; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		if ban23[x] == nil {
			ban23[x] = make(map[int]struct{})
		}
		ban23[x][y] = struct{}{}
	}

	ban34 := make([]map[int]struct{}, n3)
	var m3 int
	fmt.Fscan(in, &m3)
	for i := 0; i < m3; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		if ban34[x] == nil {
			ban34[x] = make(map[int]struct{})
		}
		ban34[x][y] = struct{}{}
	}

	dp4 := make([]int64, n4)
	copy(dp4, d)
	dp3 := compute(c, dp4, ban34)
	dp2 := compute(b, dp3, ban23)
	dp1 := compute(a, dp2, ban12)

	ans := INF
	for _, v := range dp1 {
		if v < ans {
			ans = v
		}
	}
	if ans >= INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, ans)
	}
}
