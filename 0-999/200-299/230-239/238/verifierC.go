package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n int, edges [][2]int) string {
	if n == 0 {
		return "0"
	}
	m := len(edges)
	best := m + 5
	for mask := 0; mask < (1 << m); mask++ {
		adj := make([][]int, n)
		for i, e := range edges {
			a := e[0] - 1
			b := e[1] - 1
			if (mask>>uint(i))&1 == 0 {
				adj[a] = append(adj[a], b)
			} else {
				adj[b] = append(adj[b], a)
			}
		}
		ok := false
		for s1 := 0; s1 < n && !ok; s1++ {
			for s2 := s1; s2 < n && !ok; s2++ {
				vis := make([]bool, n)
				var dfs func(int)
				dfs = func(u int) {
					if vis[u] {
						return
					}
					vis[u] = true
					for _, v := range adj[u] {
						dfs(v)
					}
				}
				dfs(s1)
				dfs(s2)
				all := true
				for i := 0; i < n; i++ {
					if !vis[i] {
						all = false
						break
					}
				}
				if all {
					ok = true
				}
			}
		}
		if ok {
			flips := bits.OnesCount(uint(mask))
			if flips < best {
				best = flips
			}
		}
	}
	return fmt.Sprintf("%d", best)
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
