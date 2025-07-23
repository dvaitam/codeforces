package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func addMod(a, b int64) int64 {
	s := a + b
	if s >= mod {
		s -= mod
	}
	return s
}

func mulMod(a, b int64) int64 {
	return (a * b) % mod
}

// polyGeom computes 1/(1-f(x)) modulo x^{k+1}, assuming f[0]==0
func polyGeom(f []int64, k int) []int64 {
	res := make([]int64, k+1)
	res[0] = 1
	for d := 1; d <= k; d++ {
		var sum int64
		for i := 1; i <= d; i++ {
			sum += f[i] * res[d-i]
			sum %= mod
		}
		res[d] = sum
	}
	return res
}

func computeClosedWalks(n, k int, edges [][2]int) []int64 {
	type Edge struct {
		to  int
		rev int
	}
	adj := make([][]Edge, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		ai := len(adj[a])
		bi := len(adj[b])
		adj[a] = append(adj[a], Edge{to: b, rev: bi})
		adj[b] = append(adj[b], Edge{to: a, rev: ai})
	}
	// allocate F[v][idx]
	F := make([][][]int64, n)
	for v := 0; v < n; v++ {
		F[v] = make([][]int64, len(adj[v]))
		for i := range F[v] {
			F[v][i] = make([]int64, k+1)
		}
	}
	parent := make([]int, n)
	parentIdx := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	root := 0
	stack := []int{root}
	order := []int{}
	parent[root] = -1
	parentIdx[root] = -1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for idx, e := range adj[v] {
			if e.to == parent[v] {
				continue
			}
			parent[e.to] = v
			parentIdx[e.to] = e.rev
			stack = append(stack, e.to)
		}
	}
	// bottom-up
	zero := make([]int64, k+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if parent[v] == -1 {
			continue
		}
		sum := make([]int64, k+1)
		for _, e := range adj[v] {
			if e.to == parent[v] {
				continue
			}
			// add shift2 of F[e.to][e.rev]
			for t := 0; t+2 <= k; t++ {
				sum[t+2] = addMod(sum[t+2], F[e.to][e.rev][t])
			}
		}
		F[v][parentIdx[v]] = polyGeom(sum, k)
	}
	// final closed walks array
	cw := make([]int64, k+1)
	var dfs func(v, p int)
	dfs = func(v, p int) {
		deg := len(adj[v])
		prefix := make([][]int64, deg+1)
		suffix := make([][]int64, deg+1)
		for i := 0; i <= deg; i++ {
			prefix[i] = make([]int64, k+1)
			suffix[i] = make([]int64, k+1)
		}
		for i := 0; i < deg; i++ {
			copy(prefix[i+1], prefix[i])
			for t := 0; t+2 <= k; t++ {
				prefix[i+1][t+2] = addMod(prefix[i+1][t+2], F[adj[v][i].to][adj[v][i].rev][t])
			}
		}
		for i := deg - 1; i >= 0; i-- {
			copy(suffix[i], suffix[i+1])
			for t := 0; t+2 <= k; t++ {
				suffix[i][t+2] = addMod(suffix[i][t+2], F[adj[v][i].to][adj[v][i].rev][t])
			}
		}
		total := prefix[deg]
		P := polyGeom(total, k)
		for t := 0; t <= k; t++ {
			cw[t] = addMod(cw[t], P[t])
		}
		for idx, e := range adj[v] {
			sum := make([]int64, k+1)
			for t := 0; t <= k; t++ {
				sum[t] = addMod(prefix[idx][t], suffix[idx+1][t])
			}
			F[v][idx] = polyGeom(sum, k)
		}
		for idx, e := range adj[v] {
			if e.to == p {
				continue
			}
			dfs(e.to, v)
		}
	}
	dfs(root, -1)
	return cw
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n1, n2, k int
	if _, err := fmt.Fscan(in, &n1, &n2, &k); err != nil {
		return
	}
	edges1 := make([][2]int, n1-1)
	for i := 0; i < n1-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		edges1[i] = [2]int{a, b}
	}
	edges2 := make([][2]int, n2-1)
	for i := 0; i < n2-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		edges2[i] = [2]int{a, b}
	}
	cw1 := computeClosedWalks(n1, k, edges1)
	cw2 := computeClosedWalks(n2, k, edges2)
	// binomial coefficients
	C := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		C[i] = make([]int64, i+1)
		C[i][0] = 1
		C[i][i] = 1
		for j := 1; j < i; j++ {
			C[i][j] = addMod(C[i-1][j-1], C[i-1][j])
		}
	}
	var ans int64
	for i := 0; i <= k; i++ {
		ans = (ans + C[k][i]*cw1[i]%mod*cw2[k-i]) % mod
	}
	fmt.Println(ans)
}
