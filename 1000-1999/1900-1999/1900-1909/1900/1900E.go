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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		val := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &val[i])
		}
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var v, u int
			fmt.Fscan(in, &v, &u)
			v--
			u--
			g[v] = append(g[v], u)
		}

		index := make([]int, n)
		low := make([]int, n)
		onStack := make([]bool, n)
		stack := make([]int, 0, n)
		compID := make([]int, n)
		compCnt := 0
		idx := 0
		compSize := []int{}
		compSum := []int64{}

		var dfs func(int)
		dfs = func(v int) {
			idx++
			index[v] = idx
			low[v] = idx
			stack = append(stack, v)
			onStack[v] = true
			for _, to := range g[v] {
				if index[to] == 0 {
					dfs(to)
					if low[to] < low[v] {
						low[v] = low[to]
					}
				} else if onStack[to] && index[to] < low[v] {
					low[v] = index[to]
				}
			}
			if low[v] == index[v] {
				size := 0
				sum := int64(0)
				for {
					w := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					onStack[w] = false
					compID[w] = compCnt
					size++
					sum += val[w]
					if w == v {
						break
					}
				}
				compSize = append(compSize, size)
				compSum = append(compSum, sum)
				compCnt++
			}
		}
		for i := 0; i < n; i++ {
			if index[i] == 0 {
				dfs(i)
			}
		}

		dag := make([]map[int]struct{}, compCnt)
		for i := 0; i < compCnt; i++ {
			dag[i] = make(map[int]struct{})
		}
		for v := 0; v < n; v++ {
			cv := compID[v]
			for _, u := range g[v] {
				cu := compID[u]
				if cv != cu {
					dag[cv][cu] = struct{}{}
				}
			}
		}

		adj := make([][]int, compCnt)
		indeg := make([]int, compCnt)
		for i := 0; i < compCnt; i++ {
			for nb := range dag[i] {
				adj[i] = append(adj[i], nb)
				indeg[nb]++
			}
		}

		order := make([]int, 0, compCnt)
		q := make([]int, 0, compCnt)
		for i := 0; i < compCnt; i++ {
			if indeg[i] == 0 {
				q = append(q, i)
			}
		}
		for qi := 0; qi < len(q); qi++ {
			v := q[qi]
			order = append(order, v)
			for _, u := range adj[v] {
				indeg[u]--
				if indeg[u] == 0 {
					q = append(q, u)
				}
			}
		}

		dpLen := make([]int, compCnt)
		dpSum := make([]int64, compCnt)
		for i := compCnt - 1; i >= 0; i-- {
			v := order[i]
			dpLen[v] = compSize[v]
			dpSum[v] = compSum[v]
			for _, u := range adj[v] {
				candLen := compSize[v] + dpLen[u]
				candSum := compSum[v] + dpSum[u]
				if candLen > dpLen[v] {
					dpLen[v] = candLen
					dpSum[v] = candSum
				} else if candLen == dpLen[v] && candSum < dpSum[v] {
					dpSum[v] = candSum
				}
			}
		}

		bestLen := 0
		bestSum := int64(0)
		first := true
		for i := 0; i < compCnt; i++ {
			if first || dpLen[i] > bestLen || (dpLen[i] == bestLen && dpSum[i] < bestSum) {
				bestLen = dpLen[i]
				bestSum = dpSum[i]
				first = false
			}
		}
		fmt.Fprintln(out, bestLen, bestSum)
	}
}
