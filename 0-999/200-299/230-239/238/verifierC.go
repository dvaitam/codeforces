package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Edge struct {
	to int
	wF int
	wB int
}

func solveCase(n int, edges [][2]int) string {
	if n == 0 {
		return "0"
	}
	adj := make([][]Edge, n)
	for _, e := range edges {
		a := e[0] - 1
		b := e[1] - 1
		adj[a] = append(adj[a], Edge{to: b, wF: 0, wB: 1})
		adj[b] = append(adj[b], Edge{to: a, wF: 1, wB: 0})
	}
	cost := make([]int, n)
	var dfs1 func(u, p int) int
	dfs1 = func(u, p int) int {
		sum := 0
		for _, e := range adj[u] {
			v := e.to
			if v == p {
				continue
			}
			sum += e.wF
			sum += dfs1(v, u)
		}
		return sum
	}
	var dfs2 func(u, p int)
	dfs2 = func(u, p int) {
		for _, e := range adj[u] {
			v := e.to
			if v == p {
				continue
			}
			cost[v] = cost[u] - e.wF + e.wB
			dfs2(v, u)
		}
	}
	cost[0] = dfs1(0, -1)
	dfs2(0, -1)
	prefix := make([]int, n+1)
	answer := cost[0]
	var dfsDelta func(u, p, depth, s1 int, best *int)
	dfsDelta = func(u, p, depth, s1 int, best *int) {
		d := depth
		mid := (d + 1) / 2
		delta := prefix[d] - prefix[mid]
		cur := cost[s1] + delta
		if cur < *best {
			*best = cur
		}
		for _, e := range adj[u] {
			v := e.to
			if v == p {
				continue
			}
			diff := -1
			if e.wF == 0 {
				diff = 1
			}
			prefix[d+1] = prefix[d] + diff
			dfsDelta(v, u, d+1, s1, best)
		}
	}
	for s1 := 0; s1 < n; s1++ {
		best := cost[s1]
		prefix[0] = 0
		dfsDelta(s1, -1, 0, s1, &best)
		if best < answer {
			answer = best
		}
	}
	return fmt.Sprintf("%d", answer)
}

type test struct{ input, expected string }

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	var tests []test
	fixed := []struct {
		n     int
		edges [][2]int
	}{
		{1, [][2]int{}},
		{2, [][2]int{{1, 2}}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", f.n))
		for _, e := range f.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(f.n, f.edges)})
	}
	for len(tests) < 100 {
		n := rng.Intn(6) + 2
		edges := make([][2]int, n-1)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			if rng.Intn(2) == 0 {
				edges[i-2] = [2]int{p, i}
			} else {
				edges[i-2] = [2]int{i, p}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(n, edges)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
