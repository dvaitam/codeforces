package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(a, e int64) int64 {
	a %= MOD
	if a < 0 {
		a += MOD
	}
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func inv(a int64) int64 {
	return modPow((a%MOD+MOD)%MOD, MOD-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type edge struct {
		to int
		w  int64
	}
	adj := make([][]edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		var c int64
		fmt.Fscan(in, &u, &v, &c)
		adj[u] = append(adj[u], edge{v, c})
		adj[v] = append(adj[v], edge{u, c})
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	parent[0] = -1
	parentW := make([]int64, n)
	order := make([]int, 0, n)
	stack := []int{0}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, e := range adj[u] {
			if e.to == parent[u] {
				continue
			}
			parent[e.to] = u
			parentW[e.to] = e.w
			stack = append(stack, e.to)
		}
	}

	A := make([]int64, n)
	B := make([]int64, n)

	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		if parent[u] == -1 {
			continue // root handled later
		}
		if len(adj[u]) == 1 { // leaf
			A[u] = 0
			B[u] = 0
			continue
		}
		var Asum, Bsum, wsum int64
		for _, e := range adj[u] {
			if e.to == parent[u] {
				continue
			}
			Asum = (Asum + A[e.to]) % MOD
			Bsum = (Bsum + B[e.to]) % MOD
			wsum = (wsum + e.w) % MOD
		}
		d := int64(len(adj[u]))
		RHS := (Bsum + (wsum+parentW[u])%MOD) % MOD
		den := ((d % MOD) - Asum + MOD) % MOD
		invDen := inv(den)
		A[u] = invDen
		B[u] = RHS * invDen % MOD
	}

	var Asum, Bsum, wsum int64
	for _, e := range adj[0] {
		Asum = (Asum + A[e.to]) % MOD
		Bsum = (Bsum + B[e.to]) % MOD
		wsum = (wsum + e.w) % MOD
	}
	d := int64(len(adj[0]))
	den := ((d % MOD) - Asum + MOD) % MOD
	res := (Bsum + wsum%MOD) % MOD * inv(den) % MOD
	fmt.Println(res)
}
