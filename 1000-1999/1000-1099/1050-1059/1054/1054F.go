package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct{ x, y int }
type VerSeg struct{ x, y1, y2 int }
type HorSeg struct{ y, x1, x2 int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var k int
	fmt.Fscan(reader, &k)
	pts := make([]Point, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &pts[i].x, &pts[i].y)
	}
	// vertical segments: sort by x, then y
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x != pts[j].x {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	ver := make([]VerSeg, 0)
	for i := 1; i < k; i++ {
		if pts[i].x == pts[i-1].x {
			ver = append(ver, VerSeg{pts[i].x, pts[i-1].y, pts[i].y})
		}
	}
	// horizontal segments: sort by y, then x
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].y != pts[j].y {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	hor := make([]HorSeg, 0)
	for i := 1; i < k; i++ {
		if pts[i].y == pts[i-1].y {
			hor = append(hor, HorSeg{pts[i].y, pts[i-1].x, pts[i].x})
		}
	}
	n := len(ver)
	m := len(hor)
	// build adjacency for bipartite graph
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if ver[i].y1 < hor[j].y && ver[i].y2 > hor[j].y &&
				hor[j].x1 < ver[i].x && hor[j].x2 > ver[i].x {
				adj[i] = append(adj[i], j)
			}
		}
	}
	// matching
	pairedL := make([]int, n)
	pairedR := make([]int, m)
	for i := 0; i < n; i++ {
		pairedL[i] = -1
	}
	for j := 0; j < m; j++ {
		pairedR[j] = -1
	}
	used := make([]bool, n)
	var dfs func(u int) bool
	dfs = func(u int) bool {
		used[u] = true
		for _, j := range adj[u] {
			if pairedR[j] < 0 || (!used[pairedR[j]] && dfs(pairedR[j])) {
				pairedL[u] = j
				pairedR[j] = u
				return true
			}
		}
		return false
	}
	for u := 0; u < n; u++ {
		if pairedL[u] < 0 {
			for i := range used {
				used[i] = false
			}
			dfs(u)
		}
	}
	// find maximum independent set via Konig's theorem
	visL := make([]bool, n)
	visR := make([]bool, m)
	var dfsL func(u int)
	var dfsR func(j int)
	dfsL = func(u int) {
		if visL[u] {
			return
		}
		visL[u] = true
		for _, j := range adj[u] {
			if pairedL[u] != j {
				dfsR(j)
			}
		}
	}
	dfsR = func(j int) {
		if visR[j] {
			return
		}
		visR[j] = true
		if pairedR[j] >= 0 {
			dfsL(pairedR[j])
		}
	}
	for u := 0; u < n; u++ {
		if pairedL[u] < 0 {
			dfsL(u)
		}
	}
	maxIndL := make([]bool, n)
	maxIndR := make([]bool, m)
	for i := 0; i < n; i++ {
		if visL[i] {
			maxIndL[i] = true
		}
	}
	for j := 0; j < m; j++ {
		if !visR[j] {
			maxIndR[j] = true
		}
	}
	// collect vertical answer segments
	type Pair struct{ x, y int }
	versAns := make([]VerSeg, 0)
	was := make(map[Pair]bool)
	for i := 0; i < n; i++ {
		if !maxIndL[i] {
			continue
		}
		x := ver[i].x
		ind := n
		for j := i; j < n; j++ {
			if !maxIndL[j] || ver[j].x != x {
				ind = j
				break
			}
			was[Pair{x, ver[j].y1}] = true
			was[Pair{x, ver[j].y2}] = true
		}
		// segment from y1 of first to y2 of last
		start := ver[i].y1
		end := ver[ind-1].y2
		versAns = append(versAns, VerSeg{x, start, end})
		i = ind - 1
	}
	// add single-point segments not covered
	for _, p := range pts {
		key := Pair{p.x, p.y}
		if !was[key] {
			versAns = append(versAns, VerSeg{p.x, p.y, p.y})
		}
	}
	// collect horizontal answer segments
	horsAns := make([]HorSeg, 0)
	was = make(map[Pair]bool)
	for i := 0; i < m; i++ {
		if !maxIndR[i] {
			continue
		}
		y := hor[i].y
		ind := m
		for j := i; j < m; j++ {
			if !maxIndR[j] || hor[j].y != y {
				ind = j
				break
			}
			was[Pair{hor[j].x1, y}] = true
			was[Pair{hor[j].x2, y}] = true
		}
		start := hor[i].x1
		end := hor[ind-1].x2
		horsAns = append(horsAns, HorSeg{y, start, end})
		i = ind - 1
	}
	for _, p := range pts {
		key := Pair{p.x, p.y}
		if !was[key] {
			horsAns = append(horsAns, HorSeg{p.y, p.x, p.x})
		}
	}
	// output horizontal then vertical
	fmt.Fprintln(writer, len(horsAns))
	for _, s := range horsAns {
		fmt.Fprintf(writer, "%d %d %d %d\n", s.x1, s.y, s.x2, s.y)
	}
	fmt.Fprintln(writer, len(versAns))
	for _, s := range versAns {
		fmt.Fprintf(writer, "%d %d %d %d\n", s.x, s.y1, s.x, s.y2)
	}
}
