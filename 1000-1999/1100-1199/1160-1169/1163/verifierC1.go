package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC1")
	cmd := exec.Command("go", "build", "-o", oracle, "1163C1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(44))
	var tests []string
	tests = append(tests, "2\n0 0\n1 1\n")
	tests = append(tests, "3\n0 0\n1 0\n0 1\n")
	for len(tests) < 100 {
		n := rng.Intn(8) + 2
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			x := rng.Intn(21) - 10
			y := rng.Intn(21) - 10
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	tests := generateTests()
	for i, input := range tests {
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
