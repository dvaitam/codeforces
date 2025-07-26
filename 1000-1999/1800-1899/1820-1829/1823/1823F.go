package main

import (
	"bufio"
	"fmt"
	"os"
)

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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, s, t int
	if _, err := fmt.Fscan(reader, &n, &s, &t); err != nil {
		return
	}
	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = len(g[i])
	}
	invDeg := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if deg[i] > 0 {
			invDeg[i] = modPow(int64(deg[i]), MOD-2)
		}
	}
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = -2
	}
	order := make([]int, 0, n)
	stack := []int{t}
	parent[t] = -1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		}
	}
	A := make([]int64, n+1)
	B := make([]int64, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		var sumA, sumB int64
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			sumA = (sumA + A[to]*invDeg[to]) % MOD
			sumB = (sumB + B[to]*invDeg[to]) % MOD
		}
		denom := (1 + MOD - sumA) % MOD
		invDen := modPow(denom, MOD-2)
		termParent := int64(0)
		if parent[v] != -1 && parent[v] != t {
			termParent = invDeg[parent[v]]
		}
		A[v] = termParent * invDen % MOD
		delta := int64(0)
		if v == s {
			delta = 1
		}
		B[v] = (delta + sumB) % MOD
		B[v] = B[v] * invDen % MOD
	}
	ans := make([]int64, n+1)
	queue := []int{t}
	ans[t] = B[t]
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			ans[to] = (A[to]*ans[v] + B[to]) % MOD
			queue = append(queue, to)
		}
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[i]%MOD)
	}
	writer.WriteByte('\n')
}
