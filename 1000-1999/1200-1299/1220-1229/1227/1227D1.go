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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	// sort indices by value descending, index ascending for tie
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		if values[idx[i]] == values[idx[j]] {
			return idx[i] < idx[j]
		}
		return values[idx[i]] > values[idx[j]]
	})

	var m int
	fmt.Fscan(in, &m)
	type query struct {
		k, pos, id int
	}
	qs := make([]query, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &qs[i].k, &qs[i].pos)
		qs[i].id = i
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].k < qs[j].k })
	answers := make([]int, m)

	// for easy version, simply compute prefix for each query
	// maintain first k indices sorted by original index
	// but since n <=100, we can do per query from scratch.
	for _, q := range qs {
		chosen := idx[:q.k]
		tmp := append([]int(nil), chosen...)
		sort.Ints(tmp)
		answers[q.id] = values[tmp[q.pos-1]]
	}

	for i := 0; i < m; i++ {
		fmt.Fprintln(out, answers[i])
	}
}
