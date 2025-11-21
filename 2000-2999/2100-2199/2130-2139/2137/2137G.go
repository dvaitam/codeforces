package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)
		from := make([]int, m)
		to := make([]int, m)
		outEdges := make([][]int, n)
		revEdges := make([][]int, n)
		deg := make([]int, n)
		indeg := make([]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			from[i] = u
			to[i] = v
			outEdges[u] = append(outEdges[u], i)
			revEdges[v] = append(revEdges[v], i)
			deg[u]++
			indeg[v]++
		}
		// topological order
		order := make([]int, 0, n)
		queue := make([]int, 0)
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				queue = append(queue, i)
			}
		}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			order = append(order, u)
			for _, e := range outEdges[u] {
				v := to[e]
				indeg[v]--
				if indeg[v] == 0 {
					queue = append(queue, v)
				}
			}
		}
		winC := make([]bool, n)
		winR := make([]bool, n)
		red := make([]bool, n)
		goodCnt := make([]int, n)
		edgeCond := make([]bool, m)

		for i := n - 1; i >= 0; i-- {
			u := order[i]
			if deg[u] == 0 {
				winC[u] = true
				winR[u] = true
				continue
			}
			good := 0
			bad := false
			for _, e := range outEdges[u] {
				v := to[e]
				cond := (!red[v]) && (deg[v] == 0 || winR[v])
				edgeCond[e] = cond
				if cond {
					good++
				}
				if red[v] || (deg[v] > 0 && !winC[v]) {
					bad = true
				}
			}
			goodCnt[u] = good
			winC[u] = good > 0
			winR[u] = !bad
		}

		queueC := make([]int, 0)
		queueR := make([]int, 0)
		enqueueC := func(u int) {
			if winC[u] {
				winC[u] = false
				queueC = append(queueC, u)
			}
		}
		enqueueR := func(u int) {
			if winR[u] {
				winR[u] = false
				queueR = append(queueR, u)
			}
		}
		processQueues := func() {
			for len(queueC) > 0 || len(queueR) > 0 {
				for len(queueR) > 0 {
					u := queueR[0]
					queueR = queueR[1:]
					for _, e := range revEdges[u] {
						p := from[e]
						if edgeCond[e] {
							edgeCond[e] = false
							if goodCnt[p] > 0 {
								goodCnt[p]--
								if goodCnt[p] == 0 && deg[p] > 0 && !red[p] {
									enqueueC(p)
								}
							}
						}
					}
				}
				for len(queueC) > 0 {
					u := queueC[0]
					queueC = queueC[1:]
					for _, e := range revEdges[u] {
						p := from[e]
						enqueueR(p)
					}
				}
			}
		}

		for ; q > 0; q-- {
			var typ, u int
			fmt.Fscan(in, &typ, &u)
			u--
			if typ == 1 {
				if !red[u] {
					red[u] = true
					enqueueC(u)
					enqueueR(u)
					processQueues()
				}
			} else {
				if winC[u] {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			}
		}
	}
}
