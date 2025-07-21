package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solve(n int) (int, []string) {
	a := n / 2
	b := n - a
	m := a * b
	edges := make([]string, 0, m)
	for i := 0; i < a; i++ {
		for j := a; j < n; j++ {
			edges = append(edges, fmt.Sprintf("%d %d", i+1, j+1))
		}
	}
	return m, edges
}

func generateCase(rng *rand.Rand) (string, int, []string) {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	m, edges := solve(n)
	return sb.String(), m, edges
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expM, expEdges := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		lines := strings.Split(out, "\n")
		if len(lines) < 1+expM {
			fmt.Fprintf(os.Stderr, "case %d failed: expected at least %d lines got %d\ninput:\n%s", i+1, 1+expM, len(lines), in)
			os.Exit(1)
		}
		gotM, err := strconv.Atoi(strings.TrimSpace(lines[0]))
		if err != nil || gotM != expM {
			fmt.Fprintf(os.Stderr, "case %d failed: expected m=%d got %s\ninput:\n%s", i+1, expM, lines[0], in)
			os.Exit(1)
		}
		for j := 0; j < expM; j++ {
			if strings.TrimSpace(lines[j+1]) != expEdges[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at edge %d: expected %q got %q\ninput:\n%s", i+1, j+1, expEdges[j], strings.TrimSpace(lines[j+1]), in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
