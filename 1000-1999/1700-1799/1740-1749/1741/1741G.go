package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type Queue struct{ data []int }

func (q *Queue) Push(x int) { q.data = append(q.data, x) }
func (q *Queue) Pop() int   { x := q.data[0]; q.data = q.data[1:]; return x }

func bfs(start int, adj [][]int) []int {
	n := len(adj) - 1
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	dist[start] = 0
	q := []int{start}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		dist1 := bfs(1, adj)
		var f int
		fmt.Fscan(in, &f)
		h := make([]int, f)
		for i := range h {
			fmt.Fscan(in, &h[i])
		}
		var k int
		fmt.Fscan(in, &k)
		p := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &p[i])
			p[i]--
		}
		hv := make([]int, k)
		for i := 0; i < k; i++ {
			hv[i] = h[p[i]]
		}
		dist := make([][]int, k)
		for i := 0; i < k; i++ {
			dist[i] = bfs(hv[i], adj)
		}
		// Precompute subset orders and validity independent of car friend
		subsetOrder := make([][]int, 1<<uint(k))
		subsetValid := make([]bool, 1<<uint(k))
		for mask := 1; mask < (1 << uint(k)); mask++ {
			nodesMap := make(map[int]int)
			order := []int{}
			for j := 0; j < k; j++ {
				if mask>>uint(j)&1 == 1 {
					v := hv[j]
					if _, ok := nodesMap[v]; !ok {
						nodesMap[v] = j
						order = append(order, j)
					}
				}
			}
			sort.Slice(order, func(a, b int) bool { return dist1[hv[order[a]]] < dist1[hv[order[b]]] })
			valid := true
			for i := 1; i < len(order); i++ {
				a := order[i-1]
				b := order[i]
				va := hv[a]
				vb := hv[b]
				if dist1[va] == dist1[vb] && va != vb {
					valid = false
					break
				}
				if dist1[va]+dist[a][vb] != dist1[vb] {
					valid = false
					break
				}
			}
			if valid {
				subsetValid[mask] = true
				subsetOrder[mask] = order
			}
		}
		// Prepare set of noCar friends
		noCar := make([]bool, f)
		for _, x := range p {
			noCar[x] = true
		}
		dp := make([]bool, 1<<uint(k))
		dp[0] = true
		for i := 0; i < f; i++ {
			if noCar[i] {
				continue
			}
			hi := h[i]
			friendMasks := []int{0}
			for mask := 1; mask < (1 << uint(k)); mask++ {
				if !subsetValid[mask] {
					continue
				}
				order := subsetOrder[mask]
				lastIdx := -1
				if len(order) > 0 {
					lastIdx = order[len(order)-1]
				}
				ok := true
				if lastIdx != -1 {
					if dist1[hv[lastIdx]]+dist[lastIdx][hi] != dist1[hi] {
						ok = false
					}
				} else {
					if dist1[hi] < 0 {
						ok = false
					}
				}
				if ok {
					friendMasks = append(friendMasks, mask)
				}
			}
			if len(friendMasks) == 1 {
				continue
			}
			newDP := make([]bool, 1<<uint(k))
			copy(newDP, dp)
			for prev := 0; prev < (1 << uint(k)); prev++ {
				if !dp[prev] {
					continue
				}
				for _, add := range friendMasks {
					m2 := prev | add
					if !newDP[m2] {
						newDP[m2] = true
					}
				}
			}
			dp = newDP
		}
		best := 0
		for mask := 0; mask < (1 << uint(k)); mask++ {
			if dp[mask] {
				bits := bits.OnesCount(uint(mask))
				if bits > best {
					best = bits
				}
			}
		}
		fmt.Fprintln(out, k-best)
	}
}
