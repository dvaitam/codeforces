package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct{ r, c int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([][]byte, n)
		var sr, sc int
		empty := 0
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			row := []byte(s)
			for j := 0; j < m; j++ {
				if row[j] != '#' {
					empty++
				}
				if row[j] == 'V' {
					sr = i
					sc = j
				}
			}
			grid[i] = row
		}
		total := empty
		// BFS from start
		INF := n*m + 5
		distS := make([]int, n*m)
		for i := range distS {
			distS[i] = INF
		}
		q := make([]int, 0)
		idx := func(r, c int) int { return r*m + c }
		startIdx := idx(sr, sc)
		distS[startIdx] = 0
		q = append(q, startIdx)
		dir := []int{-1, 0, 1, 0, -1}
		for head := 0; head < len(q); head++ {
			v := q[head]
			r := v / m
			c := v % m
			for k := 0; k < 4; k++ {
				nr := r + dir[k]
				nc := c + dir[k+1]
				if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] == '#' {
					continue
				}
				ni := idx(nr, nc)
				if distS[ni] == INF {
					distS[ni] = distS[v] + 1
					q = append(q, ni)
				}
			}
		}
		// collect accessible exits
		exits := make([]Point, 0)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if i != 0 && i != n-1 && j != 0 && j != m-1 {
					continue
				}
				if grid[i][j] == '#' {
					continue
				}
				if distS[idx(i, j)] == INF {
					continue
				}
				exits = append(exits, Point{i, j})
			}
		}
		k := len(exits)
		if k == 0 {
			fmt.Fprintln(writer, total-1)
			continue
		}
		if k == 1 {
			d := distS[idx(exits[0].r, exits[0].c)]
			fmt.Fprintln(writer, total-(d+1))
			continue
		}
		// multi-source BFS for two nearest exits
		dist1 := make([]int, n*m)
		dist2 := make([]int, n*m)
		id1 := make([]int, n*m)
		id2 := make([]int, n*m)
		for i := range dist1 {
			dist1[i] = INF
			dist2[i] = INF
			id1[i] = -1
			id2[i] = -1
		}
		type state struct {
			pos int
			id  int
			d   int
		}
		queue := make([]state, 0)
		for idxE, e := range exits {
			p := idx(e.r, e.c)
			dist1[p] = 0
			id1[p] = idxE
			queue = append(queue, state{p, idxE, 0})
		}
		for head := 0; head < len(queue); head++ {
			st := queue[head]
			p := st.pos
			id := st.id
			d := st.d
			// verify current distance matches stored
			if dist1[p] == d && id1[p] == id {
				// propagate from best record
			} else if dist2[p] == d && id2[p] == id {
				// propagate from second best
			} else {
				continue
			}
			r := p / m
			c := p % m
			for k := 0; k < 4; k++ {
				nr := r + dir[k]
				nc := c + dir[k+1]
				if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] == '#' {
					continue
				}
				np := idx(nr, nc)
				nd := d + 1
				if nd < dist1[np] {
					if id1[np] != id {
						if dist1[np] < dist2[np] {
							dist2[np] = dist1[np]
							id2[np] = id1[np]
						}
					}
					dist1[np] = nd
					id1[np] = id
					queue = append(queue, state{np, id, nd})
				} else if id1[np] != id && nd < dist2[np] {
					dist2[np] = nd
					id2[np] = id
					queue = append(queue, state{np, id, nd})
				}
			}
		}
		best := INF
		for p := 0; p < n*m; p++ {
			if distS[p] == INF {
				continue
			}
			if id1[p] == -1 || id2[p] == -1 {
				continue
			}
			if id1[p] == id2[p] {
				continue
			}
			cost := distS[p] + dist1[p] + dist2[p] + 1
			if cost < best {
				best = cost
			}
		}
		fmt.Fprintln(writer, total-best)
	}
}
