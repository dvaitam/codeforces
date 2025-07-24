package main

import (
	"bufio"
	"fmt"
	"os"
)

var n, k int
var edges [][2]int
var dist [][]int
var total int64
var positionsInit []int

func dfs(idx int, vertices []int, positions []int) {
	if idx == len(edges) {
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				total += int64(dist[positions[i]][positions[j]])
			}
		}
		return
	}
	u := edges[idx][0]
	v := edges[idx][1]
	// orientation u -> v
	vert1 := make([]int, len(vertices))
	copy(vert1, vertices)
	pos1 := make([]int, len(positions))
	copy(pos1, positions)
	if vert1[u] != -1 && vert1[v] == -1 {
		b := vert1[u]
		vert1[u] = -1
		vert1[v] = b
		pos1[b] = v
	}
	dfs(idx+1, vert1, pos1)

	// orientation v -> u
	vert2 := make([]int, len(vertices))
	copy(vert2, vertices)
	pos2 := make([]int, len(positions))
	copy(pos2, positions)
	if vert2[v] != -1 && vert2[u] == -1 {
		b := vert2[v]
		vert2[v] = -1
		vert2[u] = b
		pos2[b] = u
	}
	dfs(idx+1, vert2, pos2)
}

func bfs(start int, g [][]int) []int {
	d := make([]int, n+1)
	for i := 1; i <= n; i++ {
		d[i] = -1
	}
	q := make([]int, 0)
	q = append(q, start)
	d[start] = 0
	for head := 0; head < len(q); head++ {
		x := q[head]
		for _, y := range g[x] {
			if d[y] == -1 {
				d[y] = d[x] + 1
				q = append(q, y)
			}
		}
	}
	return d
}

func modPow(a, b, mod int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &k)
	positionsInit = make([]int, k)
	vertices := make([]int, n+1)
	for i := 0; i <= n; i++ {
		vertices[i] = -1
	}
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &positionsInit[i])
		vertices[positionsInit[i]] = i
	}
	edges = make([][2]int, n-1)
	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		u := edges[i][0]
		v := edges[i][1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	if n > 20 {
		fmt.Println(0)
		return
	}
	dist = make([][]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = bfs(i, g)
	}
	pos := make([]int, k)
	copy(pos, positionsInit)
	dfs(0, vertices, pos)
	mod := int64(998244353)
	denom := modPow(2, int64(n-1), mod)
	choose := int64(k*(k-1)/2) % mod
	denom = denom * choose % mod
	inv := modPow(denom, mod-2, mod)
	ans := total % mod * inv % mod
	fmt.Println(ans)
}
