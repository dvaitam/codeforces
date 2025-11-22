package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	u, v int
	w    int
}

type pairData struct {
	qIdx []int // indices of queries involving this pair, sorted by k
	pos  int   // first unanswered query
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)

		adj := make([][]int, n)
		edges := make([]edge, m)
		minW, maxW := int(1e9), 0
		for i := 0; i < m; i++ {
			var u, v, w int
			fmt.Fscan(in, &u, &v, &w)
			u--
			v--
			edges[i] = edge{u: u, v: v, w: w}
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			if w < minW {
				minW = w
			}
			if w > maxW {
				maxW = w
			}
		}

		type query struct {
			a, b, k int
		}
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &queries[i].a, &queries[i].b, &queries[i].k)
			queries[i].a--
			queries[i].b--
		}

		// Group queries by unordered vertex pair.
		pairs := make(map[int]*pairData)
		keyFunc := func(a, b int) int {
			if a > b {
				a, b = b, a
			}
			return a*n + b
		}
		for i := 0; i < q; i++ {
			a, b := queries[i].a, queries[i].b
			key := keyFunc(a, b)
			pd, ok := pairs[key]
			if !ok {
				pd = &pairData{}
				pairs[key] = pd
			}
			pd.qIdx = append(pd.qIdx, i)
		}
		for _, pd := range pairs {
			sort.Slice(pd.qIdx, func(i, j int) bool {
				return queries[pd.qIdx[i]].k < queries[pd.qIdx[j]].k
			})
		}

		// Initial all-edges-heavy distances: unweighted shortest path lengths.
		const inf = int(1e9)
		dist := make([][]int, n)
		for i := range dist {
			dist[i] = make([]int, n)
			for j := range dist[i] {
				dist[i][j] = inf
			}
		}
		queue := make([]int, n)
		for s := 0; s < n; s++ {
			for i := 0; i < n; i++ {
				dist[s][i] = inf
			}
			head, tail := 0, 0
			queue[tail] = s
			tail++
			dist[s][s] = 0
			for head < tail {
				u := queue[head]
				head++
				for _, v := range adj[u] {
					if dist[s][v] == inf {
						dist[s][v] = dist[s][u] + 1
						queue[tail] = v
						tail++
					}
				}
			}
		}

		// lastDist stores current distance between vertex pairs.
		lastDist := make([]int, n*n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				lastDist[i*n+j] = dist[i][j]
			}
		}

		ans := make([]int, q)
		for i := range ans {
			ans[i] = -1
		}

		// Answer queries already satisfied at smallest weight.
		for key, pd := range pairs {
			u := key / n
			v := key % n
			d := lastDist[u*n+v]
			for pd.pos < len(pd.qIdx) {
				idx := pd.qIdx[pd.pos]
				if queries[idx].k-1 >= d {
					ans[idx] = minW
					pd.pos++
				} else {
					break
				}
			}
			pairs[key] = pd
		}

		// Sort edges by weight ascending for processing.
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].w < edges[j].w
		})

		compOf := make([]int, n)
		active := make([]bool, n)
		for i := 0; i < n; i++ {
			compOf[i] = i
			active[i] = true
		}

		parent := make([]int, n)
		var find func(int) int
		find = func(x int) int {
			if parent[x] == x {
				return x
			}
			parent[x] = find(parent[x])
			return parent[x]
		}
		union := func(a, b int) {
			ra, rb := find(a), find(b)
			if ra != rb {
				parent[ra] = rb
			}
		}

		mergeComp := func(a, b int) {
			if a == b || !active[a] || !active[b] {
				return
			}
			for i := 0; i < n; i++ {
				if active[i] {
					if dist[a][i] > dist[b][i] {
						dist[a][i] = dist[b][i]
					}
					if dist[i][a] > dist[i][b] {
						dist[i][a] = dist[i][b]
					}
				}
			}
			dist[a][a] = 0
			for v := 0; v < n; v++ {
				if compOf[v] == b {
					compOf[v] = a
				}
			}
			active[b] = false
			// Floyd-style relaxation with a as intermediate.
			for i := 0; i < n; i++ {
				if !active[i] {
					continue
				}
				ia := dist[i][a]
				for j := 0; j < n; j++ {
					if !active[j] {
						continue
					}
					if ia+dist[a][j] < dist[i][j] {
						dist[i][j] = ia + dist[a][j]
					}
				}
			}
		}

		i := 0
		for i < m {
			j := i
			w := edges[i].w
			for j < m && edges[j].w == w {
				j++
			}

			for idx := 0; idx < n; idx++ {
				if active[idx] {
					parent[idx] = idx
				} else {
					parent[idx] = -1
				}
			}

			for k := i; k < j; k++ {
				u := edges[k].u
				v := edges[k].v
				cu := compOf[u]
				cv := compOf[v]
				if cu == cv {
					continue
				}
				union(cu, cv)
			}

			groups := make(map[int][]int)
			for idx := 0; idx < n; idx++ {
				if !active[idx] {
					continue
				}
				r := find(idx)
				if r == -1 {
					continue
				}
				groups[r] = append(groups[r], idx)
			}

			merged := false
			for _, comps := range groups {
				if len(comps) <= 1 {
					continue
				}
				base := comps[0]
				for _, c := range comps[1:] {
					mergeComp(base, c)
				}
				merged = true
			}

			if merged {
				for u := 0; u < n; u++ {
					for v := u + 1; v < n; v++ {
						d := dist[compOf[u]][compOf[v]]
						idx1 := u*n + v
						if d < lastDist[idx1] {
							lastDist[idx1] = d
							if pd, ok := pairs[idx1]; ok {
								for pd.pos < len(pd.qIdx) {
									qIdx := pd.qIdx[pd.pos]
									if queries[qIdx].k-1 >= d {
										if ans[qIdx] == -1 || ans[qIdx] > w {
											ans[qIdx] = w
										}
										pd.pos++
									} else {
										break
									}
								}
							}
						}
					}
				}
			}

			i = j
		}

		for i := 0; i < q; i++ {
			if ans[i] == -1 {
				ans[i] = maxW
			}
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
