package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	id int
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
		g := make([][]Edge, n)
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			y--
			edges[i] = [2]int{x, y}
			g[x] = append(g[x], Edge{y, i})
			g[y] = append(g[y], Edge{x, i})
		}

		deg := make([]int, n)
		maxDeg1, maxDeg2 := 0, 0
		cntMax1 := 0
		for i := 0; i < n; i++ {
			d := len(g[i])
			deg[i] = d
			if d > maxDeg1 {
				maxDeg2 = maxDeg1
				maxDeg1 = d
				cntMax1 = 1
			} else if d == maxDeg1 {
				cntMax1++
			} else if d > maxDeg2 {
				maxDeg2 = d
			}
		}

		bestR := int(1e9)
		root := 0
		for i := 0; i < n; i++ {
			candidate1 := (deg[i] + 1) / 2
			var maxExcl int
			if deg[i] == maxDeg1 && cntMax1 == 1 {
				maxExcl = maxDeg2
			} else {
				maxExcl = maxDeg1
			}
			candidate2 := 0
			if maxExcl > 0 {
				candidate2 = maxExcl - 1
			}
			r := candidate1
			if candidate2 > r {
				r = candidate2
			}
			if r < bestR {
				bestR = r
				root = i
			}
		}

		r := bestR
		colorForEdge := make([]int, n-1)
		nextColor := 1

		// color edges incident to root
		rootDeg := len(g[root])
		if rootDeg > 0 {
			color1 := nextColor
			nextColor++
			limit1 := r
			if limit1 > rootDeg {
				limit1 = rootDeg
			}
			for idx := 0; idx < limit1; idx++ {
				e := g[root][idx]
				colorForEdge[e.id] = color1
				dfsAssign(e.to, root, color1, &nextColor, r, g, colorForEdge)
			}
			if limit1 < rootDeg {
				color2 := nextColor
				nextColor++
				for idx := limit1; idx < rootDeg; idx++ {
					e := g[root][idx]
					colorForEdge[e.id] = color2
					dfsAssign(e.to, root, color2, &nextColor, r, g, colorForEdge)
				}
			}
		}

		fmt.Fprintln(out, r)
		for i := 0; i < n-1; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, colorForEdge[i])
		}
		out.WriteByte('\n')
	}
}

func dfsAssign(v, p, parentColor int, nextColor *int, r int, g [][]Edge, colorForEdge []int) {
	childCount := len(g[v]) - 1
	if childCount <= 0 {
		return
	}
	c := *nextColor
	*nextColor++
	for _, e := range g[v] {
		if e.to == p {
			continue
		}
		colorForEdge[e.id] = c
		dfsAssign(e.to, v, c, nextColor, r, g, colorForEdge)
	}
}
