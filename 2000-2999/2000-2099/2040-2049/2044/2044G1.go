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
		var n int
		fmt.Fscan(in, &n)

		r := make([]int, n)
		indeg := make([]int, n)
		rev := make([][]int, n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			x--
			r[i] = x
			indeg[x]++
			rev[x] = append(rev[x], i)
		}

		queue := make([]int, 0)
		head := 0
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				queue = append(queue, i)
			}
		}

		indegCopy := make([]int, n)
		copy(indegCopy, indeg)

		for head < len(queue) {
			u := queue[head]
			head++
			v := r[u]
			indegCopy[v]--
			if indegCopy[v] == 0 {
				queue = append(queue, v)
			}
		}

		depth := make([]int, n)
		for i := range depth {
			depth[i] = -1
		}

		queue = queue[:0]
		head = 0
		for i := 0; i < n; i++ {
			if indegCopy[i] > 0 {
				depth[i] = 0
				queue = append(queue, i)
			}
		}

		for head < len(queue) {
			v := queue[head]
			head++
			for _, u := range rev[v] {
				if depth[u] == -1 {
					depth[u] = depth[v] + 1
					queue = append(queue, u)
				}
			}
		}

		maxDepth := 0
		for _, d := range depth {
			if d > maxDepth {
				maxDepth = d
			}
		}

		fmt.Fprintln(out, maxDepth+2)
	}
}
