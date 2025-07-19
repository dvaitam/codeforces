package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	v, x int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, s int
	if _, err := fmt.Fscan(in, &n, &s); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	// coordinate compression
	b := make([]int, n)
	copy(b, a)
	sort.Ints(b)
	for i := 0; i < n; i++ {
		a[i] = sort.SearchInts(b, a[i])
	}
	// sorted target
	b2 := make([]int, n)
	copy(b2, a)
	sort.Ints(b2)

	// build graph of required moves
	g := make([][]edge, n)
	for i := 0; i < n; i++ {
		if a[i] != b2[i] {
			g[a[i]] = append(g[a[i]], edge{b2[i], i})
		}
	}
	cured := make([]int, n)
	var eu []int
	var rec func(int)
	rec = func(u int) {
		for cured[u] < len(g[u]) {
			e := g[u][cured[u]]
			cured[u]++
			rec(e.v)
			eu = append(eu, e.x)
		}
	}

	var cycles [][]int
	tot := 0
	for i := 0; i < n; i++ {
		eu = eu[:0]
		rec(i)
		if len(eu) > 0 {
			// copy to avoid reuse of slice
			tmp := make([]int, len(eu))
			copy(tmp, eu)
			cycles = append(cycles, tmp)
			tot += len(tmp)
		}
	}
	if tot > s {
		fmt.Fprintln(out, -1)
		return
	}
	spare := s - tot
	if spare > 2 && len(cycles) > 2 {
		if spare > len(cycles) {
			spare = len(cycles)
		}
		var small, big []int
		for i := 0; i < spare; i++ {
			last := cycles[len(cycles)-1]
			small = append(small, last[len(last)-1])
			big = append(big, last...)
			cycles = cycles[:len(cycles)-1]
		}
		// reverse small
		for i, j := 0, len(small)-1; i < j; i, j = i+1, j-1 {
			small[i], small[j] = small[j], small[i]
		}
		cycles = append(cycles, small)
		cycles = append(cycles, big)
	}
	// output
	fmt.Fprintln(out, len(cycles))
	for _, cyc := range cycles {
		fmt.Fprintln(out, len(cyc))
		for i, v := range cyc {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(out, "%d", v+1)
		}
		out.WriteByte('\n')
	}
}
