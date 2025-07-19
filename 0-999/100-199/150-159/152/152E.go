package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = 1 << 60

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m, K int
	fmt.Fscan(reader, &n, &m, &K)
	nm := n * m
	mat := make([]int, nm)
	for i := 0; i < nm; i++ {
		fmt.Fscan(reader, &mat[i])
	}
	st := make([]int, nm)
	spec := make([]int, K)
	var a, b int
	for i := 0; i < K; i++ {
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		idx := a*m + b
		st[idx] = 1 << i
		spec[i] = idx
	}
	maxMask := 1 << K
	dp := make([][]int, nm)
	preIdx := make([][]int, nm)
	preMask := make([][]int, nm)
	inq := make([][]bool, nm)
	for i := 0; i < nm; i++ {
		dp[i] = make([]int, maxMask)
		preIdx[i] = make([]int, maxMask)
		preMask[i] = make([]int, maxMask)
		inq[i] = make([]bool, maxMask)
		for s := 0; s < maxMask; s++ {
			dp[i][s] = INF
			preIdx[i][s] = -1
		}
	}
	for i := 0; i < K; i++ {
		idx := spec[i]
		mask := 1 << i
		dp[idx][mask] = mat[idx]
	}
	dx := []int{0, 1, 0, -1}
	dy := []int{1, 0, -1, 0}
	type Node struct{ u, s int }
	var queue []Node
	fullMask := maxMask - 1
	lastIdx := spec[K-1]
	for s := 1; s <= fullMask; s++ {
		queue = queue[:0]
		for u := 0; u < nm; u++ {
			if st[u] != 0 && (st[u]&s) == 0 {
				continue
			}
			for p := (s - 1) & s; p > 0; p = (p - 1) & s {
				if dp[u][p] >= INF || dp[u][s-p] >= INF {
					continue
				}
				s1 := p | st[u]
				s2 := (s - p) | st[u]
				cost := dp[u][s1] + dp[u][s2] - mat[u]
				if cost < dp[u][s] {
					dp[u][s] = cost
					preIdx[u][s] = u
					preMask[u][s] = s1
				}
			}
			if dp[u][s] < INF {
				inq[u][s] = true
				queue = append(queue, Node{u, s})
			}
		}
		for head := 0; head < len(queue); head++ {
			cur := queue[head]
			u, s0 := cur.u, cur.s
			inq[u][s0] = false
			ux := u / m
			uy := u % m
			baseCost := dp[u][s0]
			for d := 0; d < 4; d++ {
				vx := ux + dx[d]
				vy := uy + dy[d]
				if vx < 0 || vx >= n || vy < 0 || vy >= m {
					continue
				}
				v := vx*m + vy
				ts := s0 | st[v]
				nc := baseCost + mat[v]
				if nc < dp[v][ts] {
					dp[v][ts] = nc
					preIdx[v][ts] = u
					preMask[v][ts] = s0
					if !inq[v][ts] && s0 == ts {
						inq[v][ts] = true
						queue = append(queue, Node{v, ts})
					}
				}
			}
		}
	}
	res := dp[lastIdx][fullMask]
	fmt.Fprintln(writer, res)
	vis := make([]bool, nm)
	var dfs func(u, s int)
	dfs = func(u, s int) {
		if vis[u] {
			return
		}
		vis[u] = true
		pi := preIdx[u][s]
		pm := preMask[u][s]
		if pi < 0 {
			return
		}
		dfs(pi, pm)
		if pi == u {
			ns := (s - pm) | st[u]
			dfs(pi, ns)
		}
	}
	dfs(lastIdx, fullMask)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if vis[i*m+j] {
				writer.WriteByte('X')
			} else {
				writer.WriteByte('.')
			}
		}
		writer.WriteByte('\n')
	}
}
