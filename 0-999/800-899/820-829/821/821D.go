package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

type Edge struct {
	to, w int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	cells := make([][2]int, k)
	cellIndex := make(map[[2]int]int)
	rowCells := make(map[int][]int)
	colCells := make(map[int][]int)

	for i := 0; i < k; i++ {
		var r, c int
		fmt.Fscan(reader, &r, &c)
		cells[i] = [2]int{r, c}
		cellIndex[[2]int{r, c}] = i
		rowCells[r] = append(rowCells[r], i)
		colCells[c] = append(colCells[c], i)
	}

	dest := [2]int{n, m}
	destIndex, destLit := cellIndex[dest]

	rowID := make(map[int]int)
	colID := make(map[int]int)
	nextID := k

	for r := range rowCells {
		rowID[r] = nextID
		nextID++
	}
	for c := range colCells {
		colID[c] = nextID
		nextID++
	}

	if !destLit {
		destIndex = nextID
		nextID++
	}

	edges := make([][]Edge, nextID)
	add := func(u, v, w int) { edges[u] = append(edges[u], Edge{v, w}) }

	// edges for lit cells
	for i, rc := range cells {
		r, c := rc[0], rc[1]
		// adjacency
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if j, ok := cellIndex[[2]int{nr, nc}]; ok {
				add(i, j, 0)
			}
		}
		if id, ok := rowID[r]; ok {
			add(i, id, 0)
		}
		if id, ok := colID[c]; ok {
			add(i, id, 0)
		}
	}

	// edges for rows
	for r, lst := range rowCells {
		rid := rowID[r]
		for _, idx := range lst {
			add(rid, idx, 1)
		}
		if !destLit && r == dest[0] {
			add(rid, destIndex, 1)
		}
	}

	// edges for columns
	for c, lst := range colCells {
		cid := colID[c]
		for _, idx := range lst {
			add(cid, idx, 1)
		}
		if !destLit && c == dest[1] {
			add(cid, destIndex, 1)
		}
	}

	const INF = int(1e9)
	dist := make([]int, nextID)
	for i := range dist {
		dist[i] = INF
	}

	start := cellIndex[[2]int{1, 1}]
	dq := list.New()
	dist[start] = 0
	dq.PushFront(start)

	for dq.Len() > 0 {
		e := dq.Front()
		dq.Remove(e)
		u := e.Value.(int)
		du := dist[u]
		for _, ed := range edges[u] {
			nd := du + ed.w
			if nd < dist[ed.to] {
				dist[ed.to] = nd
				if ed.w == 0 {
					dq.PushFront(ed.to)
				} else {
					dq.PushBack(ed.to)
				}
			}
		}
	}

	ans := dist[destIndex]
	if ans == INF {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}
