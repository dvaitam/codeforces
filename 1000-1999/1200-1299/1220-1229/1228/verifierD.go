package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func completeTripartite(a, b, c int, r *rand.Rand) string {
	n := a + b + c
	groups := []int{}
	for i := 0; i < a; i++ {
		groups = append(groups, 1)
	}
	for i := 0; i < b; i++ {
		groups = append(groups, 2)
	}
	for i := 0; i < c; i++ {
		groups = append(groups, 3)
	}
	// shuffle vertices
	r.Shuffle(len(groups), func(i, j int) { groups[i], groups[j] = groups[j], groups[i] })
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i + 1
	}
	r.Shuffle(len(idx), func(i, j int) { idx[i], idx[j] = idx[j], idx[i] })
	edges := make([][2]int, 0, a*b+b*c+c*a)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if groups[i] != groups[j] {
				edges = append(edges, [2]int{idx[i], idx[j]})
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func randomGraph(r *rand.Rand) string {
	n := r.Intn(6) + 3
	maxEdges := n * (n - 1) / 2
	m := r.Intn(maxEdges + 1)
	edges := make(map[[2]int]struct{})
	for len(edges) < m {
		a := r.Intn(n) + 1
		b := r.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		edges[[2]int{a, b}] = struct{}{}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func generateTests() []string {
	r := rand.New(rand.NewSource(4))
	tests := make([]string, 0, 100)
	// first 10 tests complete tripartite
	for i := 0; i < 10; i++ {
		a := r.Intn(3) + 1
		b := r.Intn(3) + 1
		c := r.Intn(3) + 1
		tests = append(tests, completeTripartite(a, b, c, r))
	}
	for len(tests) < 100 {
		tests = append(tests, randomGraph(r))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	official := "./officialD"
	if err := exec.Command("go", "build", "-o", official, "1228D.go").Run(); err != nil {
		fmt.Println("failed to build official solution:", err)
		os.Exit(1)
	}
	defer os.Remove(official)
	tests := generateTests()
	for i, tc := range tests {
		exp, eerr := runBinary(official, tc)
		got, gerr := runBinary(cand, tc)
		if eerr != nil {
			fmt.Printf("official solution failed on test %d: %v\n", i+1, eerr)
			os.Exit(1)
		}
		if gerr != nil {
			fmt.Printf("candidate failed on test %d: %v\n", i+1, gerr)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
