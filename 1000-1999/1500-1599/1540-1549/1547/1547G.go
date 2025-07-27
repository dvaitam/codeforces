package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		g := make([][]int, n+1)
		rg := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			g[a] = append(g[a], b)
			rg[b] = append(rg[b], a)
		}

		reachable := make([]bool, n+1)
		stack := []int{1}
		reachable[1] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, to := range g[v] {
				if !reachable[to] {
					reachable[to] = true
					stack = append(stack, to)
				}
			}
		}

		order := make([]int, 0, n)
		used := make([]bool, n+1)
		var dfs1 func(int)
		dfs1 = func(v int) {
			used[v] = true
			for _, to := range g[v] {
				if reachable[to] && !used[to] {
					dfs1(to)
				}
			}
			order = append(order, v)
		}
		for v := 1; v <= n; v++ {
			if reachable[v] && !used[v] {
				dfs1(v)
			}
		}

		comp := make([]int, n+1)
		compCnt := 0
		var dfs2 func(int, int)
		dfs2 = func(v, c int) {
			comp[v] = c
			for _, to := range rg[v] {
				if reachable[to] && comp[to] == 0 {
					dfs2(to, c)
				}
			}
		}
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			if comp[v] == 0 {
				compCnt++
				dfs2(v, compCnt)
			}
		}

		compSize := make([]int, compCnt+1)
		for v := 1; v <= n; v++ {
			if reachable[v] {
				compSize[comp[v]]++
			}
		}

		cyc := make([]bool, compCnt+1)
		for c := 1; c <= compCnt; c++ {
			if compSize[c] > 1 {
				cyc[c] = true
			}
		}
		for v := 1; v <= n; v++ {
			if reachable[v] {
				cv := comp[v]
				for _, to := range g[v] {
					if reachable[to] && comp[to] == cv && v == to {
						cyc[cv] = true
					}
				}
			}
		}

		adjC := make([][]int, compCnt+1)
		for v := 1; v <= n; v++ {
			if !reachable[v] {
				continue
			}
			cv := comp[v]
			for _, to := range g[v] {
				if !reachable[to] {
					continue
				}
				ct := comp[to]
				if cv != ct {
					adjC[cv] = append(adjC[cv], ct)
				}
			}
		}

		startComp := comp[1]
		reachComp := make([]bool, compCnt+1)
		queue := []int{startComp}
		reachComp[startComp] = true
		for head := 0; head < len(queue); head++ {
			c := queue[head]
			for _, to := range adjC[c] {
				if !reachComp[to] {
					reachComp[to] = true
					queue = append(queue, to)
				}
			}
		}

		inf := make([]bool, compCnt+1)
		q := make([]int, 0)
		for c := 1; c <= compCnt; c++ {
			if reachComp[c] && cyc[c] {
				inf[c] = true
				q = append(q, c)
			}
		}
		for head := 0; head < len(q); head++ {
			c := q[head]
			for _, to := range adjC[c] {
				if reachComp[to] && !inf[to] {
					inf[to] = true
					q = append(q, to)
				}
			}
		}

		dp := make([]int, compCnt+1)
		indeg := make([]int, compCnt+1)
		for c := 1; c <= compCnt; c++ {
			if !reachComp[c] || inf[c] {
				continue
			}
			for _, to := range adjC[c] {
				if reachComp[to] && !inf[to] {
					indeg[to]++
				}
			}
		}

		queue = queue[:0]
		for c := 1; c <= compCnt; c++ {
			if !reachComp[c] || inf[c] {
				continue
			}
			if indeg[c] == 0 {
				queue = append(queue, c)
			}
		}
		dp[startComp] = 1
		for head := 0; head < len(queue); head++ {
			c := queue[head]
			for _, to := range adjC[c] {
				if !reachComp[to] || inf[to] {
					continue
				}
				val := dp[to] + dp[c]
				if val > 2 {
					val = 2
				}
				if val > dp[to] {
					dp[to] = val
				}
				indeg[to]--
				if indeg[to] == 0 {
					queue = append(queue, to)
				}
			}
		}

		ans := make([]int, n+1)
		for v := 1; v <= n; v++ {
			if !reachable[v] {
				ans[v] = 0
			} else if inf[comp[v]] {
				ans[v] = -1
			} else {
				ans[v] = dp[comp[v]]
			}
		}

		for v := 1; v <= n; v++ {
			if v > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[v])
		}
		fmt.Fprintln(writer)
	}
}
