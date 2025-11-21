package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u int
	v int
	w int64
}

type testCase struct {
	n     int
	edges []edge
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2,
			edges: []edge{
				{u: 1, v: 2, w: 5},
			},
		},
		{
			n: 3,
			edges: []edge{
				{1, 2, 6},
				{2, 3, 9},
			},
		},
		{
			n: 4,
			edges: []edge{
				{1, 2, 12},
				{1, 3, 6},
				{1, 4, 4},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 2
	edges := make([]edge, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		w := int64(rng.Intn(30) + 1)
		edges = append(edges, edge{u: u, v: v, w: w})
	}
	return testCase{n: n, edges: edges}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(1))
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	out = strings.TrimSpace(out)
	val, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer output %q", out)
	}
	return val, nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func bruteForce(tc testCase) int64 {
	n := tc.n
	adj := make([][]edge, n+1)
	for _, e := range tc.edges {
		adj[e.u] = append(adj[e.u], edge{v: e.v, w: e.w})
		adj[e.v] = append(adj[e.v], edge{v: e.u, w: e.w})
	}
	var best int64
	type state struct {
		node int
		g    int64
		len  int
	}
	var dfs func(int, int, int64, int)
	target := 0
	dfs = func(u, parent int, g int64, length int) {
		if g > 0 {
			val := g * int64(length)
			if val > best {
				best = val
			}
		}
		if u == target {
			if g == 0 {
				// when source==target, len=0,g=0 => product 0 handled
			}
		}
		for _, e := range adj[u] {
			if e.v == parent {
				continue
			}
			newG := g
			if newG == 0 {
				newG = e.w
			} else {
				newG = gcd(newG, e.w)
			}
			dfs(e.v, u, newG, length+1)
		}
	}
	for src := 1; src <= n; src++ {
		target = src
		dfs(src, 0, 0, 0)
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		out, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\noutput:\n%s\n", idx+1, err, out)
			os.Exit(1)
		}
		expect := bruteForce(tc)
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\ninput:\n%s\n", idx+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
