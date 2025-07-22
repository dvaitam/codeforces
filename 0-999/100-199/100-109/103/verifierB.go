package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input  string
	expect string
}

func isCthulhu(n int, edges [][2]int) bool {
	if n < 3 || len(edges) != n {
		return false
	}
	g := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	q := make([]int, 0, n)
	seen := make([]bool, n)
	q = append(q, 0)
	seen[0] = true
	for i := 0; i < len(q); i++ {
		v := q[i]
		for _, to := range g[v] {
			if !seen[to] {
				seen[to] = true
				q = append(q, to)
			}
		}
	}
	for _, s := range seen {
		if !s {
			return false
		}
	}
	return true
}

func buildCase(n int, edges [][2]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	expect := "NO\n"
	if isCthulhu(n, edges) {
		expect = "FHTAGN!\n"
	}
	return testCase{input: sb.String(), expect: expect}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	var edges [][2]int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Intn(3) == 0 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	return buildCase(n, edges)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	if strings.TrimSpace(got) != strings.TrimSpace(tc.expect) {
		return fmt.Errorf("expected %q got %q", tc.expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// simple
	cases = append(cases, buildCase(1, nil))
	cases = append(cases, buildCase(3, [][2]int{{0, 1}, {1, 2}, {2, 0}}))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
