package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Interval struct {
	high int64
	low  int64
}

func find(parent map[int64]int64, x int64) int64 {
	if x <= 0 {
		return 0
	}
	root := x
	for {
		p, ok := parent[root]
		if !ok {
			parent[root] = root
			break
		}
		if p == root {
			break
		}
		root = p
	}
	for x != root {
		p := parent[x]
		parent[x] = root
		x = p
	}
	return root
}

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
		intervals := make([]Interval, n)
		for i := 0; i < n; i++ {
			intervals[i] = Interval{high: a[i] + int64(i+1), low: a[i] + 1}
		}
		sort.Slice(intervals, func(i, j int) bool {
			if intervals[i].low == intervals[j].low {
				return intervals[i].high > intervals[j].high
			}
			return intervals[i].low > intervals[j].low
		})
		parent := make(map[int64]int64, n*2)
		res := make([]int64, 0, n)
		for _, it := range intervals {
			x := find(parent, it.high)
			if x >= it.low {
				parent[x] = find(parent, x-1)
				res = append(res, x)
			}
		}
		sort.Slice(res, func(i, j int) bool { return res[i] > res[j] })
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
