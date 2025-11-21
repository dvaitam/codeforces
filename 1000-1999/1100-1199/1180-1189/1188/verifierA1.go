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

type tree struct {
	n   int
	edg [][2]int
}

func generateTree(rng *rand.Rand, n int) tree {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{parent, i})
	}
	return tree{n: n, edg: edges}
}

func solveTree(t tree) string {
	deg := make([]int, t.n+1)
	for _, e := range t.edg {
		deg[e[0]]++
		deg[e[1]]++
	}
	for i := 1; i <= t.n; i++ {
		if deg[i] == 2 {
			return "NO"
		}
	}
	return "YES"
}

func inputString(t tree) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for _, e := range t.edg {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.ToUpper(got) != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	if bin == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Deterministic small cases
	for n := 2; n <= 6; n++ {
		for mask := 0; mask < (1 << (n - 1)); mask++ {
			t := tree{n: n}
			t.edg = make([][2]int, 0, n-1)
			parent := make([]int, n+1)
			for i := 2; i <= n; i++ {
				p := ((mask >> (i - 2)) % (i - 1)) + 1
				parent[i] = p
				t.edg = append(t.edg, [2]int{p, i})
			}
			input := inputString(t)
			exp := solveTree(t)
			if err := runCase(bin, input, exp); err != nil {
				fmt.Fprintf(os.Stderr, "deterministic case failed (n=%d): %v\ninput:\n%s", n, err, input)
				os.Exit(1)
			}
		}
	}
	// Random larger cases
	for i := 0; i < 200; i++ {
		n := rng.Intn(1000-2) + 2
		t := generateTree(rng, n)
		input := inputString(t)
		exp := solveTree(t)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "random case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
