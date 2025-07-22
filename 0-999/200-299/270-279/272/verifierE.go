package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	r := strings.NewReader(strings.TrimSpace(input))
	var n, m int
	fmt.Fscan(r, &n, &m)
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(r, &a, &b)
		a--
		b--
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	col := make([]byte, n)
	same := make([]int, n)
	q := []int{}
	for i := 0; i < n; i++ {
		same[i] = len(adj[i])
		if same[i] >= 2 {
			q = append(q, i)
		}
	}
	for qi := 0; qi < len(q); qi++ {
		v := q[qi]
		if same[v] < 2 {
			continue
		}
		old := col[v]
		col[v] ^= 1
		k := same[v]
		d := len(adj[v])
		same[v] = d - k
		for _, u := range adj[v] {
			if col[u] == old {
				same[u]--
			}
			if col[u] == col[v] {
				same[u]++
			}
			if same[u] == 2 {
				q = append(q, u)
			}
		}
	}
	for i := 0; i < n; i++ {
		if same[i] >= 2 {
			return "-1\n"
		}
	}
	out := make([]byte, n+1)
	for i := 0; i < n; i++ {
		out[i] = '0' + col[i]
	}
	out[n] = '\n'
	return string(out)
}

func generateTests() []test {
	rand.Seed(46)
	var tests []test
	fixed := []struct {
		n     int
		edges [][2]int
	}{
		{1, nil},
		{2, [][2]int{{1, 2}}},
		{3, [][2]int{{1, 2}, {2, 3}}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", f.n, len(f.edges)))
		for _, e := range f.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		n := rand.Intn(6) + 1
		deg := make([]int, n)
		edges := [][2]int{}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if deg[i] >= 3 || deg[j] >= 3 {
					continue
				}
				if rand.Intn(3) == 0 {
					edges = append(edges, [2]int{i + 1, j + 1})
					deg[i]++
					deg[j]++
				}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
