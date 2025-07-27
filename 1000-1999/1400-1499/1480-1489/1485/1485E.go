package main

import (
	"bufio"
	"fmt"
	"os"
)

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
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

		g := make([][]int, n+1)
		parent := make([]int, n+1)
		for i := 2; i <= n; i++ {
			var v int
			fmt.Fscan(in, &v)
			g[i] = append(g[i], v)
			g[v] = append(g[v], i)
		}

		val := make([]int64, n+1)
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &val[i])
		}

		depth := make([]int, n+1)
		byDepth := [][]int{{1}}
		queue := []int{1}
		parent[1] = 0
		depth[1] = 0
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range g[u] {
				if v == parent[u] {
					continue
				}
				parent[v] = u
				depth[v] = depth[u] + 1
				if len(byDepth) <= depth[v] {
					byDepth = append(byDepth, []int{})
				}
				byDepth[depth[v]] = append(byDepth[depth[v]], v)
				queue = append(queue, v)
			}
		}

		maxDepth := len(byDepth) - 1
		dp := make([]int64, n+1)
		maxVal := make([]int64, maxDepth+1)
		minVal := make([]int64, maxDepth+1)
		m1 := make([]int64, maxDepth+1)
		m2 := make([]int64, maxDepth+1)

		const inf int64 = 1<<63 - 1
		for d := maxDepth; d >= 1; d-- {
			mv := int64(-1 << 60)
			mn := int64(1 << 60)
			mm1 := int64(-1 << 60)
			mm2 := int64(-1 << 60)
			for _, v := range byDepth[d] {
				x := val[v]
				if x > mv {
					mv = x
				}
				if x < mn {
					mn = x
				}
				if dp[v]+x > mm1 {
					mm1 = dp[v] + x
				}
				if dp[v]-x > mm2 {
					mm2 = dp[v] - x
				}
			}
			maxVal[d] = mv
			minVal[d] = mn
			m1[d] = mm1
			m2[d] = mm2

			for _, u := range byDepth[d-1] {
				best := int64(0)
				for _, w := range g[u] {
					if parent[w] != u {
						continue
					}
					x := val[w]
					diff := maxVal[d] - x
					if x-minVal[d] > diff {
						diff = x - minVal[d]
					}
					c1 := dp[w] + diff

					c2 := m1[d] - x
					tmp := m2[d] + x
					if tmp > c2 {
						c2 = tmp
					}

					if c1 > best {
						best = c1
					}
					if c2 > best {
						best = c2
					}
				}
				if best > dp[u] {
					dp[u] = best
				}
			}
		}

		fmt.Fprintln(out, dp[1])
	}
}
