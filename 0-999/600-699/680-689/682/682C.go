package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	g := make([][]Edge, n+1)
	for i := 2; i <= n; i++ {
		var p int
		var c int64
		fmt.Fscan(in, &p, &c)
		g[p] = append(g[p], Edge{to: i, w: c})
	}

	type node struct {
		v       int
		dist    int64
		minPref int64
		rem     bool
	}

	stack := []node{{v: 1, dist: 0, minPref: 0, rem: false}}
	removed := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if cur.rem {
			removed++
			for _, e := range g[cur.v] {
				nd := cur.dist + e.w
				mp := cur.minPref
				if nd < mp {
					mp = nd
				}
				stack = append(stack, node{v: e.to, dist: nd, minPref: mp, rem: true})
			}
			continue
		}

		if cur.dist-cur.minPref > a[cur.v] {
			removed++
			for _, e := range g[cur.v] {
				nd := cur.dist + e.w
				mp := cur.minPref
				if nd < mp {
					mp = nd
				}
				stack = append(stack, node{v: e.to, dist: nd, minPref: mp, rem: true})
			}
			continue
		}

		for _, e := range g[cur.v] {
			nd := cur.dist + e.w
			mp := cur.minPref
			if nd < mp {
				mp = nd
			}
			stack = append(stack, node{v: e.to, dist: nd, minPref: mp, rem: false})
		}
	}

	fmt.Fprintln(out, removed)
}
