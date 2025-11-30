package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n     int
	edges [][2]int
}

// Embedded testcases (previously in testcasesE.txt) to keep verifier self contained.
const rawTestcasesE = `
4 2 4 1 2 3 4
5 2 4 3 4 3 5 1 5
6 1 5 4 6 2 3 5 6 3 5
5 2 3 2 4 2 5 1 3
6 1 2 5 6 3 6 1 6 2 5
6 3 4 1 5 1 4 5 6 2 5
3 2 3 1 3
4 2 3 1 3 3 4
3 2 3 1 2
4 2 4 1 3 1 4
6 1 5 4 5 2 6 1 6 1 3
3 1 2 1 3
6 4 6 1 4 2 3 3 6 2 5
5 1 2 4 5 2 5 1 4
5 1 2 2 4 2 5 4 5
5 2 4 2 5 1 4 1 5
6 1 5 5 6 1 6 2 5 1 3
4 2 3 1 2 1 3
4 1 2 1 3 3 4
4 2 3 3 4 1 4
4 2 3 1 3 3 4
5 2 3 3 5 1 4 1 5
6 4 5 5 6 3 6 1 6 2 5
5 2 3 2 5 1 3 3 4
4 2 3 3 4 1 4
5 1 2 2 5 3 4 1 4
6 2 4 3 4 4 6 1 3 3 5
3 2 3 1 2
4 2 3 1 3 3 4
4 2 3 2 4 1 2
4 2 3 1 2 1 3
3 2 3 1 3
6 1 2 1 5 4 6 3 6 3 5
4 1 2 3 4 1 4
3 2 3 1 2
5 4 5 1 2 3 4 3 5
6 2 4 1 2 1 5 1 6 1 3
5 2 3 1 2 2 5 3 4
5 2 4 1 2 3 4 1 5
3 1 2 1 3
4 1 2 1 3 3 4
4 2 3 3 4 1 4
4 1 3 3 4 1 4
5 1 3 3 4 3 5 1 5
4 2 3 1 3 1 4
5 2 3 4 5 1 3 1 5
5 4 5 1 2 3 4 3 5
6 1 5 1 4 2 3 2 6 3 6
6 1 5 4 5 5 6 1 3 3 5
4 2 3 3 4 1 4
3 2 3 1 3
3 1 2 1 3
6 1 4 2 3 4 5 2 6 1 3
3 2 3 1 2
5 2 3 4 5 1 3 1 5
4 1 2 3 4 1 4
3 2 3 1 3
4 2 3 2 4 1 3
5 2 4 1 2 2 5 3 5
3 1 2 1 3
5 1 2 1 3 3 5 1 5
5 2 3 2 4 1 3 3 5
4 1 2 1 3 3 4
6 2 4 3 4 1 4 2 3 1 3
5 2 3 2 4 2 5 4 5
3 2 3 1 2
3 2 3 1 2
4 2 3 2 4 1 4
4 2 3 2 4 1 3
3 2 3 1 3
3 2 3 1 2
4 1 2 1 3 3 4
4 2 3 2 4 3 4
6 1 2 3 4 4 6 2 3 5 6
5 4 5 3 4 1 4 2 4
4 2 3 1 3 3 4
4 2 3 1 2 1 3
5 1 3 3 5 1 4 1 5
6 4 6 1 4 4 5 3 6 2 5
3 2 3 1 3
3 2 3 1 2
6 2 4 1 2 4 5 5 6 3 6
3 1 2 1 3
4 2 3 2 4 1 3
3 2 3 1 2
4 2 3 1 3 3 4
3 2 3 1 2
4 2 3 2 4 1 2
5 2 4 2 5 3 4 1 4
6 1 2 3 4 4 5 5 6 3 5
5 2 4 1 4 1 3 3 5
3 1 2 1 3
5 1 3 3 4 1 4 1 5
3 1 2 1 3
3 2 3 1 3
5 4 5 2 4 1 2 1 5
4 2 3 2 4 1 2
3 2 3 1 3
6 1 2 3 4 4 6 4 5 5 6
6 3 4 2 3 3 6 2 5 1 3
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcasesE, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		expectedTokens := 1 + 2*(n-1)
		if len(fields) != expectedTokens {
			return nil, fmt.Errorf("line %d: expected %d tokens got %d", idx+1, expectedTokens, len(fields))
		}
		edges := make([][2]int, 0, n-1)
		for i := 1; i < len(fields); i += 2 {
			a, err := strconv.Atoi(fields[i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge a: %w", idx+1, err)
			}
			b, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge b: %w", idx+1, err)
			}
			edges = append(edges, [2]int{a, b})
		}
		cases = append(cases, testCase{n: n, edges: edges})
	}
	return cases, nil
}

type bitset struct {
	data []uint64
}

func newBitset(n int) *bitset {
	size := (n + 64) / 64
	return &bitset{data: make([]uint64, size)}
}

func (b *bitset) set(i int) {
	b.data[i/64] |= 1 << (uint(i) % 64)
}

func (b *bitset) get(i int) bool {
	return (b.data[i/64] & (1 << (uint(i) % 64))) != 0
}

func (b *bitset) shiftOr(s int) {
	if s <= 0 {
		return
	}
	wordShift := s / 64
	bitShift := uint(s % 64)
	n := len(b.data)
	for i := n - 1; i >= 0; i-- {
		var v uint64
		j := i - wordShift
		if j >= 0 {
			v = b.data[j] << bitShift
			if bitShift > 0 && j-1 >= 0 {
				v |= b.data[j-1] >> (64 - bitShift)
			}
		}
		b.data[i] |= v
	}
}

// solve211ECase mirrors 211E.go to compute expected output.
func solve211ECase(tc testCase) string {
	n := tc.n
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	parent := make([]int, n+1)
	siz := make([]int, n+1)
	vis := make([]bool, n+1)
	var dfs func(u, p int)
	dfs = func(u, p int) {
		if vis[u] {
			return
		}
		vis[u] = true
		parent[u] = p
		siz[u] = 1
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs(v, u)
			siz[u] += siz[v]
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			dfs(i, 0)
		}
	}

	possible := make([]bool, n+1)
	for v := 1; v <= n; v++ {
		deg := len(adj[v])
		if deg < 2 {
			continue
		}
		sizes := make([]int, 0, deg)
		for _, u := range adj[v] {
			if parent[u] == v {
				sizes = append(sizes, siz[u])
			} else {
				sizes = append(sizes, n-siz[v])
			}
		}
		bs := newBitset(n)
		bs.set(0)
		for _, s := range sizes {
			bs.shiftOr(s)
		}
		total := n - 1
		for a := 1; a <= total-1; a++ {
			if bs.get(a) {
				possible[a] = true
			}
		}
	}

	total := n - 1
	var ans [][2]int
	for a := 1; a <= total-1; a++ {
		if possible[a] {
			ans = append(ans, [2]int{a, total - a})
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(len(ans)))
	for _, p := range ans {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		expected := solve211ECase(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(tc.n))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf(" %d %d", e[0], e[1]))
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d failed: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
