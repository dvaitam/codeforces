package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"sort"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)

	type Point struct{ r, c int }
	cells := make([]Point, k)
	cellMap := make(map[Point]int)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &cells[i].r, &cells[i].c)
		cellMap[cells[i]] = i
	}

	comp := make([]int, k)
	for i := range comp {
		comp[i] = -1
	}
	compCount := 0
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for i := 0; i < k; i++ {
		if comp[i] == -1 {
			q := []int{i}
			comp[i] = compCount
			head := 0
			for head < len(q) {
				u := q[head]
				head++
				for j := 0; j < 4; j++ {
					nr, nc := cells[u].r+dr[j], cells[u].c+dc[j]
					if vIdx, ok := cellMap[Point{nr, nc}]; ok && comp[vIdx] == -1 {
						comp[vIdx] = compCount
						q = append(q, vIdx)
					}
				}
			}
			compCount++
		}
	}

	adj := make([][]int, compCount)
	compR := make([][]int, compCount)
	compC := make([][]int, compCount)
	for i := 0; i < k; i++ {
		compR[comp[i]] = append(compR[comp[i]], cells[i].r)
		compC[comp[i]] = append(compC[comp[i]], cells[i].c)
	}

	for i := 0; i < compCount; i++ {
		sort.Ints(compR[i])
		sort.Ints(compC[i])
	}

	for i := 0; i < compCount; i++ {
		for j := i + 1; j < compCount; j++ {
			can_bridge := false
			// Check rows
			for _, r1 := range compR[i] {
				idx := sort.SearchInts(compR[j], r1-2)
				for l := idx; l < len(compR[j]) && compR[j][l] <= r1+2; l++ {
					can_bridge = true
					break
				}
				if can_bridge { break }
			}
			if can_bridge {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
				continue
			}
			// Check cols
			for _, c1 := range compC[i] {
				idx := sort.SearchInts(compC[j], c1-2)
				for l := idx; l < len(compC[j]) && compC[j][l] <= c1+2; l++ {
					can_bridge = true
					break
				}
				if can_bridge { break }
			}
			if can_bridge {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	const INF = 1 << 30
	dist := make([]int, compCount)
	for i := range dist {
		dist[i] = INF
	}

	startIdx, startIsLit := cellMap[Point{1, 1}]
	if !startIsLit {
		fmt.Fprintln(writer, -1)
		return
	}
	startComp := comp[startIdx]
	dist[startComp] = 0

	dq := list.New()
	dq.PushFront(startComp)

	for dq.Len() > 0 {
		e := dq.Front()
		dq.Remove(e)
		u := e.Value.(int)

		for _, v := range adj[u] {
			if dist[v] > dist[u]+1 {
				dist[v] = dist[u] + 1
				dq.PushBack(v)
			}
		}
	}

	ans := INF
	targetIdx, targetIsLit := cellMap[Point{n, m}]
	if targetIsLit {
		ans = min(ans, dist[comp[targetIdx]])
	}

	for i := 0; i < k; i++ {
		r, c := cells[i].r, cells[i].c
		if dist[comp[i]] != INF {
			if abs(r-n) <= 1 || abs(c-m) <= 1 {
				ans = min(ans, dist[comp[i]]+1)
			}
		}
	}

	if ans == INF {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}