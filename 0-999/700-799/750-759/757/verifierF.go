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

func generateCase(rng *rand.Rand) (string, string, error) {
	n := rng.Intn(5) + 2
	m := n - 1 + rng.Intn(3)
	s := rng.Intn(n) + 1
	type edge struct{ u, v int }
	edges := make([]edge, 0, m)
	// build tree first
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	exist := make(map[[2]int]bool)
	for _, e := range edges {
		a, b := e.u, e.v
		if a > b {
			a, b = b, a
		}
		exist[[2]int{a, b}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if exist[[2]int{a, b}] {
			continue
		}
		exist[[2]int{a, b}] = true
		edges = append(edges, edge{u, v})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, len(edges), s)
	for _, e := range edges {
		w := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, w)
	}
	input := sb.String()
	exp, err := runCase("757F.go", input)
	if err != nil {
		return "", "", err
	}
	return input, exp, nil
}

func runCase(exe, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp, err := generateCase(rng)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate case: %v", err)
			os.Exit(1)
		}
		got, err := runCase(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
