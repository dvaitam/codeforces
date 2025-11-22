package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to int
	id int
}

type pair struct {
	coef int64
	cnst int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	out := bufio.NewWriter(os.Stdout)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		var s string
		fmt.Fscan(in, &s)

		path := make([][2]int, 0, n+m-1)
		x, y := 0, 0
		path = append(path, [2]int{x, y})
		for _, ch := range s {
			if ch == 'D' {
				x++
			} else {
				y++
			}
			path = append(path, [2]int{x, y})
		}

		onPath := make([][]bool, n)
		for i := range onPath {
			onPath[i] = make([]bool, m)
		}
		for _, p := range path {
			onPath[p[0]][p[1]] = true
		}

		grid := make([][]int64, n)
		rowSum := make([]int64, n)
		colSum := make([]int64, m)
		for i := 0; i < n; i++ {
			grid[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &grid[i][j])
				if !onPath[i][j] {
					rowSum[i] += grid[i][j]
					colSum[j] += grid[i][j]
				}
			}
		}

		nodes := n + m
		edgesCnt := len(path)
		adj := make([][]edge, nodes)
		cellOfEdge := make([][2]int, edgesCnt)

		for id, p := range path {
			cellOfEdge[id] = p
			u := p[0]
			v := n + p[1]
			adj[u] = append(adj[u], edge{to: v, id: id})
			adj[v] = append(adj[v], edge{to: u, id: id})
		}

		coefArr := make([]int64, nodes)
		cnstArr := make([]int64, nodes)

		var dfsCoef func(v, parent int) pair
		dfsCoef = func(v, parent int) pair {
			var sumCoef, sumConst int64
			for _, e := range adj[v] {
				if e.to == parent {
					continue
				}
				res := dfsCoef(e.to, v)
				sumCoef += res.coef
				sumConst += res.cnst
			}
			var known int64
			if v < n {
				known = rowSum[v]
			} else {
				known = colSum[v-n]
			}

			coef := 1 - sumCoef
			cnst := -known - sumConst
			coefArr[v] = coef
			cnstArr[v] = cnst
			return pair{coef: coef, cnst: cnst}
		}

		root := 0 // row 0
		rootPair := dfsCoef(root, -1)

		var targetSum int64
		if rootPair.coef == 0 {
			targetSum = 0
		} else {
			targetSum = -rootPair.cnst / rootPair.coef
		}

		edgeVal := make([]int64, edgesCnt)

		var dfsAssign func(v, parent int) int64
		dfsAssign = func(v, parent int) int64 {
			var sumChild int64
			for _, e := range adj[v] {
				if e.to == parent {
					continue
				}
				val := dfsAssign(e.to, v)
				edgeVal[e.id] = val
				sumChild += val
			}

			var known int64
			if v < n {
				known = rowSum[v]
			} else {
				known = colSum[v-n]
			}

			need := targetSum - known - sumChild
			return need
		}

		dfsAssign(root, -1)

		for id, p := range cellOfEdge {
			grid[p[0]][p[1]] = edgeVal[id]
		}

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, grid[i][j])
			}
			if t > 1 || i+1 < n {
				fmt.Fprintln(out)
			}
		}
	}
	out.Flush()
}
