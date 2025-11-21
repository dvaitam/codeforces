package main

import (
	"bufio"
	"fmt"
	"os"
)

type nodeState struct {
	v, parent int
	parity    int
	diff      int64
	minDiff   int64
	maxDiff   int64
	idx       int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		ans := make([]int64, n+1)
		stack := []nodeState{{v: 1, parent: 0, parity: 0, diff: a[1], minDiff: 0, maxDiff: 0, idx: 0}}

		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			v := top.v
			if top.idx == 0 {
				if top.parity == 0 {
					ans[v] = top.diff - top.minDiff
				} else {
					ans[v] = top.maxDiff - top.diff
				}
			}
			if top.idx < len(adj[v]) {
				to := adj[v][top.idx]
				top.idx++
				if to == top.parent {
					continue
				}
				childParity := 1 - top.parity
				var childDiff int64
				if childParity == 0 {
					childDiff = top.diff + a[to]
				} else {
					childDiff = top.diff - a[to]
				}
				childMin := top.minDiff
				if top.diff < childMin {
					childMin = top.diff
				}
				childMax := top.maxDiff
				if top.diff > childMax {
					childMax = top.diff
				}
				stack = append(stack, nodeState{v: to, parent: v, parity: childParity, diff: childDiff, minDiff: childMin, maxDiff: childMax, idx: 0})
			} else {
				stack = stack[:len(stack)-1]
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
