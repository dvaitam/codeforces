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
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, filepath.Join(dir, "605B.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2
	maxEdges := n * (n - 1) / 2
	m := n - 1 + rng.Intn(maxEdges-(n-1)+1)
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	weight := 1
	for i := 0; i < n-1; i++ {
		weight += rng.Intn(5) + 1
		fmt.Fprintf(&b, "%d 1\n", weight)
	}
	for i := n - 1; i < m; i++ {
		weight += rng.Intn(5) + 1
		fmt.Fprintf(&b, "%d 0\n", weight)
	}
	return b.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
