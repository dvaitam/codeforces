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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	typ := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &typ[i])
	}
	adj := make([][]int, n)
	indeg := make([]int, n)
	for i := 0; i < m; i++ {
		var t1, t2 int
		fmt.Fscan(in, &t1, &t2)
		adj[t2] = append(adj[t2], t1)
		indeg[t1]++
	}

	qCPU := make([]int, 0)
	qGPU := make([]int, 0)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			if typ[i] == 0 {
				qCPU = append(qCPU, i)
			} else {
				qGPU = append(qGPU, i)
			}
		}
	}

	processed := 0
	calls := 0

	for processed < n {
		for len(qCPU) > 0 {
			u := qCPU[0]
			qCPU = qCPU[1:]
			processed++
			for _, v := range adj[u] {
				indeg[v]--
				if indeg[v] == 0 {
					if typ[v] == 0 {
						qCPU = append(qCPU, v)
					} else {
						qGPU = append(qGPU, v)
					}
				}
			}
		}
		if processed >= n {
			break
		}
		if len(qGPU) > 0 {
			calls++
			for len(qGPU) > 0 {
				u := qGPU[0]
				qGPU = qGPU[1:]
				processed++
				for _, v := range adj[u] {
					indeg[v]--
					if indeg[v] == 0 {
						if typ[v] == 0 {
							qCPU = append(qCPU, v)
						} else {
							qGPU = append(qGPU, v)
						}
					}
				}
			}
		}
	}

	fmt.Fprintln(out, calls)
}
