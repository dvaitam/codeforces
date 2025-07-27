package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

func powMod(a, b int) int {
	res := 1
	x := a % MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * x % MOD
		}
		x = x * x % MOD
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	// distances between all pairs
	dist := make([][]int, n)
	for s := 0; s < n; s++ {
		d := make([]int, n)
		for i := 0; i < n; i++ {
			d[i] = -1
		}
		q := []int{s}
		d[s] = 0
		for front := 0; front < len(q); front++ {
			v := q[front]
			for _, nb := range adj[v] {
				if d[nb] == -1 {
					d[nb] = d[v] + 1
					q = append(q, nb)
				}
			}
		}
		dist[s] = d
	}

	parent := make([][]int, n)
	sub := make([][]int, n)
	for i := 0; i < n; i++ {
		parent[i] = make([]int, n)
		sub[i] = make([]int, n)
	}

	type pair struct{ v, p int }
	dfs := func(root int) {
		stack := []pair{{root, -1}}
		order := make([]pair, 0)
		parent[root][root] = root
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			v, p := cur.v, cur.p
			parent[root][v] = p
			order = append(order, cur)
			for _, nb := range adj[v] {
				if nb == p {
					continue
				}
				stack = append(stack, pair{nb, v})
			}
		}
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i].v
			p := order[i].p
			s := 1
			for _, nb := range adj[v] {
				if nb == p {
					continue
				}
				s += sub[root][nb]
			}
			sub[root][v] = s
		}
	}

	for r := 0; r < n; r++ {
		dfs(r)
	}

	// DP table F[i][j] = probability u before v when distances are i and j
	F := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		F[i] = make([]int, n+1)
	}
	for j := 0; j <= n; j++ {
		F[0][j] = 1
	}
	inv2 := (MOD + 1) / 2
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			F[i][j] = (F[i-1][j] + F[i][j-1]) * inv2 % MOD
		}
	}
	invN := powMod(n, MOD-2)

	subtreeSize := func(root, node int) int {
		if root == node {
			return 0
		}
		x := node
		for parent[root][x] != root {
			x = parent[root][x]
		}
		return sub[root][x]
	}

	ans := 0
	for a := 0; a < n; a++ {
		for b := 0; b < a; b++ {
			prob := 0
			for w := 0; w < n; w++ {
				if dist[a][w]+dist[b][w] == dist[a][b] {
					su := subtreeSize(w, a)
					sv := subtreeSize(w, b)
					weight := (n - su - sv) % MOD
					weight = weight * invN % MOD
					pa := F[dist[w][a]][dist[w][b]]
					prob = (prob + pa*weight) % MOD
				}
			}
			ans = (ans + prob) % MOD
		}
	}

	fmt.Println(ans)
}
