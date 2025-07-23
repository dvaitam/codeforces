package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

// powmod computes a^b mod MOD
func powmod(a, b int64) int64 {
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

func modInv(a int64) int64 {
	return powmod(a, MOD-2)
}

type Edge struct {
	to    int
	w     int64
	color int // 0 red, 1 black
}

var (
	g      [][]Edge
	up     [][]int
	depth  []int
	diff   []int // red - black from root
	length []int
	prod   []int64
)

func dfs(v, p int, w int64, c int) {
	if p != 0 {
		depth[v] = depth[p] + 1
		diff[v] = diff[p]
		if c == 0 {
			diff[v]++
		} else {
			diff[v]--
		}
		length[v] = length[p] + 1
		prod[v] = prod[p] * w % MOD
	} else {
		prod[v] = 1
	}
	up[0][v] = p
	for _, e := range g[v] {
		if e.to == p {
			continue
		}
		dfs(e.to, v, e.w, e.color)
	}
}

func buildLCA(n int) {
	for k := 1; k < len(up); k++ {
		for v := 1; v <= n; v++ {
			if up[k-1][v] != 0 {
				up[k][v] = up[k-1][up[k-1][v]]
			}
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	// lift a
	diffDepth := depth[a] - depth[b]
	for k := 0; diffDepth > 0; k++ {
		if diffDepth&1 == 1 {
			a = up[k][a]
		}
		diffDepth >>= 1
	}
	if a == b {
		return a
	}
	for k := len(up) - 1; k >= 0; k-- {
		if up[k][a] != up[k][b] {
			a = up[k][a]
			b = up[k][b]
		}
	}
	return up[0][a]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	n := 0
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g = make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		var x int64
		var c int
		fmt.Fscan(in, &u, &v, &x, &c)
		g[u] = append(g[u], Edge{v, x, c})
		g[v] = append(g[v], Edge{u, x, c})
	}
	LOG := 0
	for (1 << LOG) <= n {
		LOG++
	}
	up = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	diff = make([]int, n+1)
	length = make([]int, n+1)
	prod = make([]int64, n+1)
	dfs(1, 0, 0, 0)
	buildLCA(n)

	ans := int64(1)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			l := lca(i, j)
			lenPath := length[i] + length[j] - 2*length[l]
			if lenPath == 0 {
				continue
			}
			d := diff[i] + diff[j] - 2*diff[l]
			if absInt(3*d) <= lenPath {
				// compute product of weights along path
				prodPath := prod[i] * prod[j] % MOD
				inv := modInv(prod[l])
				prodPath = prodPath * inv % MOD
				prodPath = prodPath * inv % MOD
				ans = ans * prodPath % MOD
			}
		}
	}
	fmt.Println(ans)
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
