package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	u, v int
}

const MOD int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		loops := make([]int, n+1)
		adj := make([][]int, n+1)
		edges := make([]Edge, 0, n)
		loopsCount := 0
		for i := 0; i < n; i++ {
			u, v := a[i], b[i]
			if u == v {
				loops[u]++
				loopsCount++
			} else {
				id := len(edges)
				edges = append(edges, Edge{u, v})
				adj[u] = append(adj[u], id)
				adj[v] = append(adj[v], id)
			}
		}

		valid := true
		for j := 1; j <= n && valid; j++ {
			if loops[j] > 1 {
				valid = false
			}
		}
		if valid {
			deg := make([]int, n+1)
			for j := 1; j <= n; j++ {
				deg[j] = len(adj[j])
			}
			for j := 1; j <= n && valid; j++ {
				if loops[j] == 0 && deg[j] == 0 {
					valid = false
				}
			}
			if valid {
				assigned := make([]bool, n+1)
				q := make([]int, 0)
				for j := 1; j <= n; j++ {
					if loops[j] == 1 {
						assigned[j] = true
						q = append(q, j)
					}
				}
				for j := 1; j <= n; j++ {
					if !assigned[j] && deg[j] == 1 {
						q = append(q, j)
					}
				}
				processed := make([]bool, len(edges))
				head := 0
				for head < len(q) && valid {
					j := q[head]
					head++
					if assigned[j] {
						for _, id := range adj[j] {
							if processed[id] {
								continue
							}
							u, v := edges[id].u, edges[id].v
							other := v
							if j == v {
								other = u
							}
							if assigned[other] {
								valid = false
								break
							}
							assigned[other] = true
							processed[id] = true
							deg[other]--
							deg[j]--
							q = append(q, other)
						}
						deg[j] = 0
					} else {
						if deg[j] == 0 {
							valid = false
							break
						}
						id := -1
						for _, e := range adj[j] {
							if !processed[e] {
								id = e
								break
							}
						}
						if id == -1 {
							valid = false
							break
						}
						u, v := edges[id].u, edges[id].v
						other := v
						if j == v {
							other = u
						}
						processed[id] = true
						assigned[j] = true
						deg[j]--
						deg[other]--
						if !assigned[other] && deg[other] == 1 {
							q = append(q, other)
						}
					}
				}
				if valid {
					visited := make([]bool, n+1)
					cycleCount := 0
					for v := 1; v <= n && valid; v++ {
						if !assigned[v] && deg[v] > 0 && !visited[v] {
							stack := []int{v}
							visited[v] = true
							vertices := 0
							edgesCount := 0
							for len(stack) > 0 {
								cur := stack[len(stack)-1]
								stack = stack[:len(stack)-1]
								vertices++
								for _, id := range adj[cur] {
									if processed[id] {
										continue
									}
									edgesCount++
									u, w := edges[id].u, edges[id].v
									other := w
									if cur == w {
										other = u
									}
									if !visited[other] {
										visited[other] = true
										stack = append(stack, other)
									}
								}
							}
							edgesCount /= 2
							if edgesCount != vertices {
								valid = false
								break
							}
							cycleCount++
						}
					}
					if valid {
						ans := modPow(int64(n), int64(loopsCount))
						ans = ans * modPow(2, int64(cycleCount)) % MOD
						fmt.Fprintln(out, ans)
						continue
					}
				}
			}
		}
		if !valid {
			fmt.Fprintln(out, 0)
		}
	}
}
