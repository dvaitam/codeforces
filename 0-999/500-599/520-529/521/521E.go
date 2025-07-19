package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxN = 200005

var (
	n, m         int
	h            []int
	ent          int
	tim, cnt     int
	vis, dep, fa []int
	U, V, W, R   []int
	edges        []Edge
	out          *bufio.Writer
)

type Edge struct {
	to, next int
}

func link(u, v int) {
	ent++
	edges[ent] = Edge{to: v, next: h[u]}
	h[u] = ent
}

func get(x, y int) []int {
	var a []int
	for x != y {
		a = append(a, x)
		x = fa[x]
	}
	a = append(a, y)
	return a
}

func printVec(a []int) {
	fmt.Fprintf(out, "%d", len(a))
	for _, v := range a {
		fmt.Fprintf(out, " %d", v)
	}
	fmt.Fprintln(out)
}

func printPath(st, en, x, y int) {
	a := get(x, st)
	reverse(a)
	b := get(en, y)
	reverse(b)
	a = append(a, b...)
	printVec(a)
}

func printAns(x, y int) {
	fmt.Fprintln(out, "YES")
	// find st
	st := W[cnt]
	for dep[st] != V[cnt] {
		st = fa[st]
	}
	// find en
	var en int
	if dep[y] > U[cnt] {
		en = y
	} else {
		en = R[cnt]
	}
	// first path
	a := get(st, en)
	printVec(a)
	// second and third
	printPath(st, en, x, y)
	printPath(st, en, W[cnt], R[cnt])
	out.Flush()
	os.Exit(0)
}

func dfs(o, ft, d int) {
	tim++
	vis[o] = tim
	dep[o] = d
	for x := h[o]; x != 0; x = edges[x].next {
		y := edges[x].to
		if x == (ft ^ 1) {
			continue
		}
		if vis[y] == 0 {
			fa[y] = o
			dfs(y, x, d+1)
			for cnt > 0 && U[cnt] == d {
				cnt--
			}
			if cnt > 0 && V[cnt] > d {
				V[cnt] = d
			}
		} else if vis[y] < vis[o] {
			if cnt > 0 && dep[y] < V[cnt] {
				printAns(o, y)
			}
			for cnt > 0 && U[cnt] > dep[y] {
				cnt--
			}
			cnt++
			U[cnt] = dep[y]
			V[cnt] = dep[o]
			R[cnt] = y
			W[cnt] = o
		}
	}
}

func reverse(a []int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fscan(in, &n, &m)
	sizeE := 2*m + 5
	edges = make([]Edge, sizeE)
	h = make([]int, n+5)
	vis = make([]int, n+5)
	dep = make([]int, n+5)
	fa = make([]int, n+5)
	U = make([]int, m+5)
	V = make([]int, m+5)
	W = make([]int, m+5)
	R = make([]int, m+5)
	ent = 1
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		link(u, v)
		link(v, u)
	}
	for i := 1; i <= n; i++ {
		if vis[i] == 0 {
			dfs(i, 0, 1)
		}
	}
	fmt.Fprintln(out, "NO")
}
