package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edgeKey struct{ u, v int }

func key(a, b int) edgeKey {
	if a > b {
		a, b = b, a
	}
	return edgeKey{a, b}
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-n+1) + n - 1
	if m%2 == 1 {
		if m < maxEdges {
			m++
		} else {
			m--
		}
	}
	edges := make([]edgeKey, 0, m)
	used := make(map[edgeKey]bool)
	// generate tree first
	for i := 2; i <= n; i++ {
		j := rng.Intn(i-1) + 1
		ek := key(i, j)
		edges = append(edges, ek)
		used[ek] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		ek := key(u, v)
		if used[ek] {
			continue
		}
		edges = append(edges, ek)
		used[ek] = true
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&in, "%d %d\n", e.u, e.v)
	}
	return in.String()
}

func solvePossible(n, m int) bool { return m%2 == 0 }

func checkAnswer(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return fmt.Errorf("input parse: %v", err)
	}
	edges := make(map[edgeKey]bool)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		edges[key(a, b)] = true
	}

	out := bufio.NewReader(strings.NewReader(output))
	tok, _ := out.Peek(1)
	if len(tok) == 0 {
		return fmt.Errorf("empty output")
	}
	if string(tok) == "N" {
		line, _ := out.ReadString('\n')
		line = strings.TrimSpace("N" + line)
		if line != "No solution" {
			return fmt.Errorf("expected 'No solution' got %q", line)
		}
		if solvePossible(n, m) {
			return fmt.Errorf("solution exists but got No solution")
		}
		return nil
	}
	if !solvePossible(n, m) {
		return fmt.Errorf("should output No solution")
	}
	need := m / 2
	used := make(map[edgeKey]bool)
	for i := 0; i < need; i++ {
		var x, y, z int
		if _, err := fmt.Fscan(out, &x, &y, &z); err != nil {
			return fmt.Errorf("read path %d: %v", i+1, err)
		}
		if x < 1 || x > n || y < 1 || y > n || z < 1 || z > n {
			return fmt.Errorf("vertex out of range")
		}
		e1 := key(x, y)
		e2 := key(y, z)
		if !edges[e1] || !edges[e2] {
			return fmt.Errorf("edge not found")
		}
		if used[e1] || used[e2] {
			return fmt.Errorf("edge used twice")
		}
		used[e1] = true
		used[e2] = true
	}
	// ensure all edges used
	if len(used) != len(edges) {
		return fmt.Errorf("not all edges used")
	}
	return nil
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return checkAnswer(input, buf.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
