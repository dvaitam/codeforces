package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseD struct {
	n int
	m int
	a []int
	b []int
}

func generateTests() []testCaseD {
	rand.Seed(42)
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rand.Intn(6) + 2 // 2..7
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges + 1)
		a := make([]int, m)
		b := make([]int, m)
		edges := map[[2]int]bool{}
		for j := 0; j < m; j++ {
			for {
				u := rand.Intn(n)
				v := rand.Intn(n)
				if u == v {
					continue
				}
				if u > v {
					u, v = v, u
				}
				key := [2]int{u, v}
				if !edges[key] {
					edges[key] = true
					a[j] = u
					b[j] = v
					break
				}
			}
		}
		tests[i] = testCaseD{n: n, m: m, a: a, b: b}
	}
	return tests
}

func solveD(t testCaseD) (string, string) {
	N := t.n
	M := t.m
	A := make([]int, M)
	B := make([]int, M)
	copy(A, t.a)
	copy(B, t.b)
	edges := make([][]int, N)
	for i := 0; i < N; i++ {
		edges[i] = []int{}
	}
	for i := 0; i < M; i++ {
		u, v := A[i], B[i]
		edges[u] = append(edges[u], v)
		edges[v] = append(edges[v], u)
	}
	for i := 0; i < N; i++ {
		sort.Ints(edges[i])
	}
	next := make([]int, N)
	prev := make([]int, N)
	inUnvis := make([]bool, N)
	for i := 0; i < N; i++ {
		next[i] = i + 1
		prev[i] = i - 1
		inUnvis[i] = true
	}
	if N > 0 {
		next[N-1] = -1
	}
	head := 0
	remove := func(u int) {
		inUnvis[u] = false
		if u == head {
			head = next[u]
		}
		if prev[u] != -1 {
			next[prev[u]] = next[u]
		}
		if next[u] != -1 {
			prev[next[u]] = prev[u]
		}
	}
	P := make([]int, N)
	Q := make([]int, N)
	num := 1
	var dfs func(n, u, ls int) bool
	dfs = func(n, u, ls int) bool {
		k := 0
		ok := false
		var nx []int
		var nd []int
		for x := head; x != -1; x = next[x] {
			for k < len(edges[n]) && edges[n][k] < x {
				k++
			}
			if k < len(edges[n]) && edges[n][k] == x {
				continue
			}
			nx = append(nx, x)
		}
		if ls == 1 {
			nd = append(nd, u)
		}
		for _, x := range nx {
			remove(x)
		}
		for i, v := range nx {
			if i == len(nx)-1 && len(nd) == 0 {
				dfs(v, n, 1)
				ok = true
			} else {
				if dfs(v, n, 0) {
					nd = append(nd, v)
				}
			}
		}
		if len(nd) > 0 {
			P[n] = num
			for i, v := range nd {
				P[v] = num + i + 1
			}
			Q[n] = num + len(nd)
			for i, v := range nd {
				Q[v] = num + i
			}
			num += len(nd) + 1
			ok = true
		}
		if !ok && u == -1 {
			P[n] = num
			Q[n] = num
			num++
			ok = true
		}
		return !ok
	}
	for i := 0; i < N; i++ {
		if inUnvis[i] {
			remove(i)
			dfs(i, -1, 0)
		}
	}
	var sb1, sb2 strings.Builder
	for i := 0; i < N; i++ {
		if i > 0 {
			sb1.WriteByte(' ')
		}
		fmt.Fprint(&sb1, P[i])
	}
	for i := 0; i < N; i++ {
		if i > 0 {
			sb2.WriteByte(' ')
		}
		fmt.Fprint(&sb2, Q[i])
	}
	return sb1.String(), sb2.String()
}

func buildInput(tests []testCaseD) string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tests))
	for _, t := range tests {
		fmt.Fprintf(&b, "%d %d\n", t.n, t.m)
		for i := 0; i < t.m; i++ {
			fmt.Fprintf(&b, "%d %d\n", t.a[i]+1, t.b[i]+1)
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	input := buildInput(tests)

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "execution failed:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(&out)
	outputs := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			outputs = append(outputs, line)
		}
	}
	if len(outputs) != len(tests)*2 {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", len(tests)*2, len(outputs))
		os.Exit(1)
	}
	for i, t := range tests {
		exp1, exp2 := solveD(t)
		if outputs[2*i] != exp1 || outputs[2*i+1] != exp2 {
			fmt.Fprintf(os.Stderr, "test %d failed\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
