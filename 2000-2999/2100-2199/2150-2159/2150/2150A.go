package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Disjoint-set over occupied cells to jump to the next free position.
func find(x int64, parent map[int64]int64) int64 {
	if nxt, ok := parent[x]; ok {
		root := find(nxt, parent)
		parent[x] = root
		return root
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		var s string
		fmt.Fscan(in, &s)

		parent := make(map[int64]int64, m+n)

		initPos := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &initPos[i])
		}
		for _, p := range initPos {
			if _, ok := parent[p]; !ok {
				parent[p] = find(p+1, parent)
			}
		}

		pos := int64(1)
		for i := 0; i < n; i++ {
			if s[i] == 'A' {
				pos++
			} else { // 'B'
				pos = find(pos+1, parent)
			}
			if _, ok := parent[pos]; !ok {
				parent[pos] = find(pos+1, parent)
			}
		}

		black := make([]int64, 0, len(parent))
		for k := range parent {
			black = append(black, k)
		}
		sort.Slice(black, func(i, j int) bool { return black[i] < black[j] })

		fmt.Fprintln(out, len(black))
		for i, v := range black {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
