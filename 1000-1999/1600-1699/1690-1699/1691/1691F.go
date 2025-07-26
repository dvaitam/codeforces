package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(a, e int64) int64 {
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

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(in, &n, &k)

	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modPow(fact[n], MOD-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	chooseK := make([]int64, n+1)
	if k <= n {
		for i := k; i <= n; i++ {
			chooseK[i] = fact[i] * invFact[k] % MOD * invFact[i-k] % MOD
		}
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	order := make([]int, 0, n)
	order = append(order, 0)
	for idx := 0; idx < len(order); idx++ {
		v := order[idx]
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			order = append(order, to)
		}
	}

	size := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		sz := 1
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			sz += size[to]
		}
		size[v] = sz
	}

	var result int64
	for x := 0; x < n; x++ {
		childSize := make([]int, len(adj[x]))
		var sumChildren int64
		for idx, y := range adj[x] {
			s := 0
			if parent[y] == x {
				s = size[y]
			} else {
				s = n - size[x]
			}
			childSize[idx] = s
			sumChildren = (sumChildren + chooseK[s]) % MOD
		}

		cnt := chooseK[n] - sumChildren
		if cnt < 0 {
			cnt += MOD
		}
		result = (result + int64(n)%MOD*cnt) % MOD

		for _, sComp := range childSize {
			comp := sComp
			subtree := n - comp
			cnt := chooseK[subtree] - (sumChildren - chooseK[comp])
			cnt %= MOD
			if cnt < 0 {
				cnt += MOD
			}
			add := int64(comp) * int64(subtree) % MOD
			add = add * cnt % MOD
			result = (result + add) % MOD
		}
	}

	if result < 0 {
		result += MOD
	}
	fmt.Println(result % MOD)
}
