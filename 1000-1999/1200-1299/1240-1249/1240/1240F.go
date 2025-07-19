package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	if k > n {
		k = n
	}
	// read edges
	type edge struct{ u, v int }
	edges := make([]edge, m)
	adj := make([][]int, n)
	pos := make([][2]int, m)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		edges[i] = edge{x, y}
		pos[i][0] = len(adj[x])
		adj[x] = append(adj[x], i)
		pos[i][1] = len(adj[y])
		adj[y] = append(adj[y], i)
	}
	// initialize random colors
	rand.Seed(11)
	color := make([]int, m)
	cnt := make([][]int, n)
	for i := 0; i < n; i++ {
		cnt[i] = make([]int, k+1)
	}
	for i := 0; i < m; i++ {
		c := rand.Intn(k) + 1
		color[i] = c
		u, v := edges[i].u, edges[i].v
		cnt[u][c]++
		cnt[v][c]++
	}
	// queue of nodes to process
	inq := make([]bool, n)
	queue := make([]int, 0, n)
	add := func(x int) {
		if !inq[x] && chk(cnt[x], k) {
			inq[x] = true
			queue = append(queue, x)
		}
	}
	// initial nodes
	for i := 0; i < n; i++ {
		add(i)
	}
	// process
	head := 0
	for head < len(queue) {
		x := queue[head]
		head++
		inq[x] = false
		mi, ma := minmax(cnt[x], k)
		if ma-mi <= 2 {
			continue
		}
		// work on node x
		if rand.Int()&1 == 1 {
			// pick incident edge with color count >= mi+2
			var posA []int
			for _, eid := range adj[x] {
				if cnt[x][color[eid]] >= mi+2 {
					posA = append(posA, eid)
				}
			}
			eid := posA[rand.Intn(len(posA))]
			// pick tb with cnt == mi
			var Bs []int
			for c := 1; c <= k; c++ {
				if cnt[x][c] == mi {
					Bs = append(Bs, c)
				}
			}
			tb := Bs[rand.Intn(len(Bs))]
			// modify edge eid to color tb
			u, v := edges[eid].u, edges[eid].v
			oc := color[eid]
			if oc != tb {
				cnt[u][oc]--
				cnt[v][oc]--
				color[eid] = tb
				cnt[u][tb]++
				cnt[v][tb]++
				add(u)
				add(v)
			}
		} else {
			// pick incident edge with color count == ma
			var posA []int
			for _, eid := range adj[x] {
				if cnt[x][color[eid]] == ma {
					posA = append(posA, eid)
				}
			}
			eid := posA[rand.Intn(len(posA))]
			// pick tb with cnt <= ma-2
			var Bs []int
			for c := 1; c <= k; c++ {
				if cnt[x][c] <= ma-2 {
					Bs = append(Bs, c)
				}
			}
			tb := Bs[rand.Intn(len(Bs))]
			u, v := edges[eid].u, edges[eid].v
			oc := color[eid]
			if oc != tb {
				cnt[u][oc]--
				cnt[v][oc]--
				color[eid] = tb
				cnt[u][tb]++
				cnt[v][tb]++
				add(u)
				add(v)
			}
		}
	}
	// output
	for i := 0; i < m; i++ {
		fmt.Fprintln(writer, color[i])
	}
}

// chk whether node violates
func chk(cnt []int, k int) bool {
	mi, ma := minmax(cnt, k)
	return ma-mi > 2
}

// minmax returns min and max count for colors 1..k
func minmax(cnt []int, k int) (int, int) {
	mi := cnt[1]
	ma := cnt[1]
	for c := 2; c <= k; c++ {
		if cnt[c] < mi {
			mi = cnt[c]
		}
		if cnt[c] > ma {
			ma = cnt[c]
		}
	}
	return mi, ma
}
