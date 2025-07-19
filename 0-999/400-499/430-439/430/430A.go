package main

import (
	"fmt"
	"sort"
)

func main() {
	var n, m int
	if _, err := fmt.Scan(&n, &m); err != nil {
		return
	}
	type pair struct{ val, idx int }
	a := make([]pair, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i].val)
		a[i].idx = i
	}
	sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })
	ans := make([]int, n)
	q := 0
	for _, p := range a {
		ans[p.idx] = q
		q ^= 1
	}
	for i, v := range ans {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
