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

func generateTree(n int, rng *rand.Rand) [][3]int {
	edges := make([][3]int, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i) + 1
		w := rng.Intn(10) + 1
		edges[i-1] = [3]int{p, i + 1, w}
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string, error) {
	n := rng.Intn(5) + 2
	q := rng.Intn(5) + 1
	seq := rng.Perm(n)
	for i := range seq {
		seq[i]++
	}
	edges := generateTree(n, rng)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i, v := range seq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			v := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "1 %d %d %d\n", l, r, v)
		} else {
			x := rng.Intn(n-1) + 1
			fmt.Fprintf(&sb, "2 %d\n", x)
		}
	}
	input := sb.String()
	exp, err := runCase("757G.go", input)
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
