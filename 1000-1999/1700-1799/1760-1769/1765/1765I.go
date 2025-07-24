package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Interval struct{ l, r int }

func addInterval(list *[]Interval, l, r int) {
	if r < l {
		return
	}
	*list = append(*list, Interval{l, r})
}

func mergeIntervals(list []Interval) []Interval {
	if len(list) == 0 {
		return list
	}
	sort.Slice(list, func(i, j int) bool { return list[i].l < list[j].l })
	res := []Interval{list[0]}
	for _, iv := range list[1:] {
		last := &res[len(res)-1]
		if iv.l <= last.r+1 {
			if iv.r > last.r {
				last.r = iv.r
			}
		} else {
			res = append(res, iv)
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var xs, ys, xt, yt, n int
	if _, err := fmt.Fscan(in, &xs, &ys); err != nil {
		return
	}
	fmt.Fscan(in, &xt, &yt)
	fmt.Fscan(in, &n)
	type Piece struct {
		t    byte
		x, y int
	}
	pieces := make([]Piece, n)
	rowPieces := make([][]int, 9)
	pieceMap := make(map[[2]int]bool)
	minX, maxX := xs, xs
	if xt < minX {
		minX = xt
	}
	if xt > maxX {
		maxX = xt
	}
	for i := 0; i < n; i++ {
		var t string
		fmt.Fscan(in, &t, &pieces[i].x, &pieces[i].y)
		pieces[i].t = t[0]
		rowPieces[pieces[i].y] = append(rowPieces[pieces[i].y], pieces[i].x)
		pieceMap[[2]int{pieces[i].x, pieces[i].y}] = true
		if pieces[i].x < minX {
			minX = pieces[i].x
		}
		if pieces[i].x > maxX {
			maxX = pieces[i].x
		}
	}
	for y := 1; y <= 8; y++ {
		sort.Ints(rowPieces[y])
	}
	minBound := minX - 8
	maxBound := maxX + 8

	attacked := make([][]Interval, 9)

	isPiece := func(x, y int) bool {
		return pieceMap[[2]int{x, y}]
	}

	addAttackCell := func(x, y int) {
		if y < 1 || y > 8 {
			return
		}
		if x < minBound || x > maxBound {
			return
		}
		addInterval(&attacked[y], x, x)
	}

	for _, p := range pieces {
		x, y := p.x, p.y
		addInterval(&attacked[y], x, x) // occupied square
		switch p.t {
		case 'K':
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					addAttackCell(x+dx, y+dy)
				}
			}
		case 'N':
			moves := [][2]int{{1, 2}, {2, 1}, {-1, 2}, {-2, 1}, {1, -2}, {2, -1}, {-1, -2}, {-2, -1}}
			for _, mv := range moves {
				addAttackCell(x+mv[0], y+mv[1])
			}
		case 'R', 'Q':
			arr := rowPieces[y]
			idx := sort.SearchInts(arr, x)
			left := minBound
			if idx > 0 {
				left = arr[idx-1] + 1
			}
			addInterval(&attacked[y], left, x-1)
			right := maxBound
			if idx+1 < len(arr) {
				right = arr[idx+1] - 1
			}
			addInterval(&attacked[y], x+1, right)
			fallthrough
		case 'B':
			if p.t == 'B' || p.t == 'Q' {
				dirs := [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
				for _, d := range dirs {
					nx, ny := x+d[0], y+d[1]
					for nx >= minBound && nx <= maxBound && ny >= 1 && ny <= 8 {
						if isPiece(nx, ny) {
							break
						}
						addAttackCell(nx, ny)
						nx += d[0]
						ny += d[1]
					}
				}
			}
			if p.t == 'R' {
				continue
			}
		}
		if p.t == 'R' || p.t == 'Q' {
			// vertical directions
			for dy := y + 1; dy <= 8; dy++ {
				if isPiece(x, dy) {
					break
				}
				addAttackCell(x, dy)
			}
			for dy := y - 1; dy >= 1; dy-- {
				if isPiece(x, dy) {
					break
				}
				addAttackCell(x, dy)
			}
		}
	}

	for y := 1; y <= 8; y++ {
		attacked[y] = mergeIntervals(attacked[y])
	}

	// collect candidate x coordinates
	candMap := make(map[int]bool)
	for d := -7; d <= 7; d++ {
		candMap[xs+d] = true
		candMap[xt+d] = true
	}
	for _, p := range pieces {
		for d := -7; d <= 7; d++ {
			candMap[p.x+d] = true
		}
	}
	candMap[minBound] = true
	candMap[maxBound] = true
	for y := 1; y <= 8; y++ {
		for _, iv := range attacked[y] {
			candMap[iv.l] = true
			candMap[iv.l-1] = true
			candMap[iv.r] = true
			candMap[iv.r+1] = true
		}
	}

	X := make([]int, 0, len(candMap))
	for x := range candMap {
		if x >= minBound-1 && x <= maxBound+1 {
			X = append(X, x)
		}
	}
	sort.Ints(X)
	uniq := X[:0]
	prev := 999999999999
	for _, x := range X {
		if prev != x {
			uniq = append(uniq, x)
			prev = x
		}
	}
	X = uniq

	index := make(map[int]int)
	for i, x := range X {
		index[x] = i
	}

	m := len(X)
	safe := make([][]bool, 9)
	for y := 1; y <= 8; y++ {
		safe[y] = make([]bool, m)
		j := 0
		ivs := attacked[y]
		for i, x := range X {
			for j < len(ivs) && ivs[j].r < x {
				j++
			}
			if j < len(ivs) && ivs[j].l <= x && x <= ivs[j].r {
				safe[y][i] = false
			} else {
				safe[y][i] = true
			}
		}
	}

	nodeID := func(ix, y int) int { return ix*8 + (y - 1) }
	totalNodes := m * 8
	type Edge struct{ to, cost int }
	graph := make([][]Edge, totalNodes)

	// horizontal edges
	for y := 1; y <= 8; y++ {
		prev := -1
		for i := 0; i < m; i++ {
			if !safe[y][i] {
				continue
			}
			if prev != -1 {
				// check if segment from prev to i is free
				x1, x2 := X[prev], X[i]
				if x1 > x2 {
					x1, x2 = x2, x1
				}
				ivs := attacked[y]
				j := sort.Search(len(ivs), func(j int) bool { return ivs[j].r >= x1 })
				good := true
				if j < len(ivs) && ivs[j].l <= x2 {
					good = false
				}
				if good {
					c := x2 - x1
					u := nodeID(prev, y)
					v := nodeID(i, y)
					graph[u] = append(graph[u], Edge{v, c})
					graph[v] = append(graph[v], Edge{u, c})
				} else {
					prev = i
					continue
				}
			}
			prev = i
		}
	}

	// vertical and diagonal step-1 edges
	for i, x := range X {
		for y := 1; y <= 8; y++ {
			if !safe[y][i] {
				continue
			}
			u := nodeID(i, y)
			// horizontal step by 1
			for _, dx := range []int{-1, 1} {
				nx := x + dx
				if j, ok := index[nx]; ok && safe[y][j] {
					v := nodeID(j, y)
					graph[u] = append(graph[u], Edge{v, 1})
				}
			}
			// vertical and diagonal
			for dy := -1; dy <= 1; dy++ {
				ny := y + dy
				if ny < 1 || ny > 8 {
					continue
				}
				if dy == 0 {
					continue
				}
				if safe[ny][i] {
					v := nodeID(i, ny)
					graph[u] = append(graph[u], Edge{v, 1})
				}
				for _, dx := range []int{-1, 1} {
					nx := x + dx
					if j, ok := index[nx]; ok && ny >= 1 && ny <= 8 && safe[ny][j] {
						v := nodeID(j, ny)
						graph[u] = append(graph[u], Edge{v, 1})
					}
				}
			}
		}
	}

	startIdx := index[xs]
	targetIdx := index[xt]
	startNode := nodeID(startIdx, ys)
	targetNode := nodeID(targetIdx, yt)

	const INF int64 = 1 << 60
	dist := make([]int64, totalNodes)
	for i := range dist {
		dist[i] = INF
	}
	dist[startNode] = 0

	type Item struct {
		node int
		dist int64
	}
	pq := make([]Item, 0)
	push := func(it Item) {
		pq = append(pq, it)
		i := len(pq) - 1
		for i > 0 {
			p := (i - 1) / 2
			if pq[p].dist <= pq[i].dist {
				break
			}
			pq[p], pq[i] = pq[i], pq[p]
			i = p
		}
	}
	pop := func() Item {
		it := pq[0]
		last := pq[len(pq)-1]
		pq = pq[:len(pq)-1]
		if len(pq) > 0 {
			pq[0] = last
			i := 0
			for {
				l := i*2 + 1
				if l >= len(pq) {
					break
				}
				r := l + 1
				c := l
				if r < len(pq) && pq[r].dist < pq[l].dist {
					c = r
				}
				if pq[i].dist <= pq[c].dist {
					break
				}
				pq[i], pq[c] = pq[c], pq[i]
				i = c
			}
		}
		return it
	}

	push(Item{startNode, 0})
	for len(pq) > 0 {
		it := pop()
		if it.dist != dist[it.node] {
			continue
		}
		if it.node == targetNode {
			fmt.Println(it.dist)
			return
		}
		for _, e := range graph[it.node] {
			nd := it.dist + int64(e.cost)
			if nd < dist[e.to] {
				dist[e.to] = nd
				push(Item{e.to, nd})
			}
		}
	}
	fmt.Println(-1)
}
