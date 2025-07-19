package main

import (
	"bufio"
	"fmt"
	"os"
)

type triple struct{ a, b, c int }
type pair struct{ from, to int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, ai, bi, ci int
	fmt.Fscan(reader, &n, &ai, &bi, &ci)
	// zero-based
	ai--
	bi--
	ci--
	// read adjacency and build edges
	cor := make([][]int, n)
	for i := 0; i < n; i++ {
		cor[i] = make([]int, n)
		for j := 0; j < n; j++ {
			cor[i][j] = -1
		}
	}
	edge := make([][26][]int, n)
	// read chars one by one
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var ch byte
			// skip whitespace
			for {
				b, err := reader.ReadByte()
				if err != nil {
					return
				}
				if b != ' ' && b != '\n' && b != '\r' && b != '\t' {
					ch = b
					break
				}
			}
			if ch != '*' {
				v := int(ch - 'a')
				cor[i][j] = v
				edge[i][v] = append(edge[i][v], j)
			}
		}
	}
	// BFS state
	dist := make([][][]int, n)
	prev := make([][][]triple, n)
	last := make([][][]pair, n)
	for i := 0; i < n; i++ {
		dist[i] = make([][]int, n)
		prev[i] = make([][]triple, n)
		last[i] = make([][]pair, n)
		for j := 0; j < n; j++ {
			dist[i][j] = make([]int, n)
			prev[i][j] = make([]triple, n)
			last[i][j] = make([]pair, n)
		}
	}
	// initial state (0,1,2) and its permutations (as index orders)
	start := [3]int{0, 1, 2}
	permIdx := [][3]int{
		{0, 1, 2}, {0, 2, 1}, {1, 0, 2},
		{1, 2, 0}, {2, 0, 1}, {2, 1, 0},
	}
	// mark initial distances
	for _, pi := range permIdx {
		p0 := start[pi[0]]
		p1 := start[pi[1]]
		p2 := start[pi[2]]
		dist[p0][p1][p2] = 1
	}
	// BFS queue holds one representative ordering
	queue := []triple{{start[0], start[1], start[2]}}
	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		d := dist[cur.a][cur.b][cur.c]
		// try moving each token i
		for i := 0; i < 3; i++ {
			// indices of other two
			var oa, ob int
			if i == 0 {
				oa, ob = 1, 2
			} else if i == 1 {
				oa, ob = 0, 2
			} else {
				oa, ob = 0, 1
			}
			A := []int{cur.a, cur.b, cur.c}[oa]
			B := []int{cur.a, cur.b, cur.c}[ob]
			u := cor[A][B]
			if u < 0 {
				continue
			}
			from := []int{cur.a, cur.b, cur.c}[i]
			for _, v := range edge[from][u] {
				if v == A || v == B {
					continue
				}
				// new triple in this order
				var nt [3]int
				nt[0], nt[1], nt[2] = cur.a, cur.b, cur.c
				nt[i] = v
				// if not visited in this ordering
				if dist[nt[0]][nt[1]][nt[2]] == 0 {
					// mark all permutations of new triple
					for _, pi := range permIdx {
						p0 := nt[pi[0]]
						p1 := nt[pi[1]]
						p2 := nt[pi[2]]
						dist[p0][p1][p2] = d + 1
						prev[p0][p1][p2] = cur
						last[p0][p1][p2] = pair{from: from, to: v}
					}
					queue = append(queue, triple{nt[0], nt[1], nt[2]})
				}
			}
		}
	}
	// target state
	if dist[ai][bi][ci] == 0 {
		fmt.Println(-1)
		return
	}
	steps := dist[ai][bi][ci] - 1
	fmt.Println(steps)
	// reconstruct
	A, B, C := ai, bi, ci
	for dist[A][B][C] != 1 {
		lp := last[A][B][C]
		// print move: to, from (1-based)
		fmt.Println(lp.to+1, lp.from+1)
		prv := prev[A][B][C]
		A, B, C = prv.a, prv.b, prv.c
	}
}
