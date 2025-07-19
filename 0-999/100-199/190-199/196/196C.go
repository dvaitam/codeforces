package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	N      int
	g      [][]int
	mm     []struct{ x, y, order, idx int }
	u      []bool
	mc     []int
	L      int
	px, py int
)

func dfs(v int) int {
	u[v] = true
	r := 1
	for _, w := range g[v] {
		if !u[w] {
			r += dfs(w)
		}
	}
	mc[v] = r
	return r
}

func vect(ax, ay, bx, by int) int64 {
	return int64(ax)*int64(by) - int64(ay)*int64(bx)
}

func swap(i, j int) {
	mm[i], mm[j] = mm[j], mm[i]
}

func quickSort(l, r int) {
	i, j := l, r
	x := mm[(i+j)/2].x
	y := mm[(i+j)/2].y
	for i <= j {
		for i <= r && vect(mm[i].x-px, mm[i].y-py, x-px, y-py) < 0 {
			i++
		}
		for j >= l && vect(mm[j].x-px, mm[j].y-py, x-px, y-py) > 0 {
			j--
		}
		if i <= j {
			swap(i, j)
			i++
			j--
		}
	}
	if l < j {
		quickSort(l, j)
	}
	if i < r {
		quickSort(i, r)
	}
}

func rec(a, v int) {
	mm[a].order = v
	u[v] = true
	swap(L, a)
	px = mm[L].x
	py = mm[L].y
	// sort by angle
	if mc[v] > 1 {
		quickSort(L+1, L+mc[v]-1)
	}
	L++
	for _, w := range g[v] {
		if !u[w] {
			rec(L, w)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &N)
	g = make([][]int, N+1)
	mm = make([]struct{ x, y, order, idx int }, N)
	u = make([]bool, N+1)
	mc = make([]int, N+1)
	for i := 0; i < N-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	m := 0
	for i := 0; i < N; i++ {
		fmt.Fscan(reader, &mm[i].x, &mm[i].y)
		mm[i].order = 0
		mm[i].idx = i
		if mm[i].y > mm[m].y {
			m = i
		}
	}
	// compute subtree sizes
	dfs(1)
	// prepare for recursion
	for i := range u {
		u[i] = false
	}
	// start recursive ordering
	rec(m, 1)
	// build result
	res := make([]int, N)
	for i := 0; i < N; i++ {
		res[mm[i].idx] = mm[i].order
	}
	// output
	for i := 0; i < N; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, res[i])
	}
	writer.WriteByte('\n')
}
