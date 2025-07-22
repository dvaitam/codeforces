package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	exe := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", exe, filepath.Join(dir, "475E.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1) / 2
	extra := rng.Intn(min(5, maxEdges-(n-1)+1))
	edges := make(map[[2]int]struct{})
	for i := 2; i <= n; i++ {
		j := rng.Intn(i-1) + 1
		a, b := i, j
		if a > b {
			a, b = b, a
		}
		edges[[2]int{a, b}] = struct{}{}
	}
	for len(edges) < (n-1)+extra {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		edges[[2]int{a, b}] = struct{}{}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
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
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
