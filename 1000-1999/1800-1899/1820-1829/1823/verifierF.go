package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 998244353

type Test struct {
	n     int
	s     int
	t     int
	edges [][2]int
}

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

func solve(tc Test) []int64 {
	n, s, t := tc.n, tc.s, tc.t
	g := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
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
	res := make([]int64, n)
	for i := 1; i <= n; i++ {
		res[i-1] = ans[i] % MOD
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(7)
	const cases = 100
	tests := make([]Test, cases)
	for i := range tests {
		n := rand.Intn(6) + 2
		s := rand.Intn(n) + 1
		t := rand.Intn(n) + 1
		for t == s {
			t = rand.Intn(n) + 1
		}
		edges := make([][2]int, n-1)
		for j := 1; j < n; j++ {
			u := rand.Intn(j) + 1
			v := j + 1
			edges[j-1] = [2]int{u, v}
		}
		tests[i] = Test{n: n, s: s, t: t, edges: edges}
	}

	for idx, tc := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.s, tc.t)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error running binary on test %d: %v\n", idx+1, err)
			fmt.Print(out.String())
			return
		}
		reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
		exp := solve(tc)
		for i := 0; i < tc.n; i++ {
			var val int64
			if _, err := fmt.Fscan(reader, &val); err != nil {
				fmt.Printf("test %d: failed to read output\n", idx+1)
				return
			}
			if val%MOD != exp[i] {
				fmt.Printf("test %d: expected %d got %d at vertex %d\n", idx+1, exp[i], val%MOD, i+1)
				return
			}
		}
	}
	fmt.Printf("verified %d test cases\n", len(tests))
}
