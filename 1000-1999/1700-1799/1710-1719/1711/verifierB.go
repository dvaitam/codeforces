package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestCase struct {
	n     int
	m     int
	a     []int
	edges [][2]int
}

func generateTests() []TestCase {
	r := rand.New(rand.NewSource(42))
	tests := make([]TestCase, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(10) + 1 // 1..10
		maxEdges := n * (n - 1) / 2
		m := r.Intn(maxEdges + 1)
		a := make([]int, n+1)
		for j := 1; j <= n; j++ {
			a[j] = r.Intn(10)
		}
		edges := make([][2]int, 0, m)
		used := make(map[[2]int]bool)
		for len(edges) < m {
			u := r.Intn(n) + 1
			v := r.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			p := [2]int{u, v}
			if used[p] {
				continue
			}
			used[p] = true
			edges = append(edges, p)
		}
		tests[i] = TestCase{n: n, m: m, a: a, edges: edges}
	}
	return tests
}

func computeAnswer(tc TestCase) int {
	if tc.m%2 == 0 {
		return 0
	}
	deg := make([]int, tc.n+1)
	for _, e := range tc.edges {
		deg[e[0]]++
		deg[e[1]]++
	}
	ans := int(1e9)
	for i := 1; i <= tc.n; i++ {
		if deg[i]%2 == 1 && tc.a[i] < ans {
			ans = tc.a[i]
		}
	}
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		if deg[u]%2 == 0 && deg[v]%2 == 0 {
			if tc.a[u]+tc.a[v] < ans {
				ans = tc.a[u] + tc.a[v]
			}
		}
	}
	return ans
}

func buildInput(tests []TestCase) string {
	var b bytes.Buffer
	fmt.Fprintln(&b, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for i := 1; i <= tc.n; i++ {
			fmt.Fprint(&b, tc.a[i])
			if i < tc.n {
				fmt.Fprint(&b, " ")
			}
		}
		fmt.Fprintln(&b)
		for _, e := range tc.edges {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	input := buildInput(tests)
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "binary execution failed:", err)
		os.Exit(1)
	}
	tokens := strings.Fields(out.String())
	if len(tokens) != len(tests) {
		fmt.Fprintln(os.Stderr, "wrong number of outputs:", len(tokens))
		os.Exit(1)
	}
	for i, tok := range tokens {
		got, err := strconv.Atoi(tok)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid integer output: %s\n", i+1, tok)
			os.Exit(1)
		}
		expect := computeAnswer(tests[i])
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d: expected %d, got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
