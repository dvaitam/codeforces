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

type testG struct {
	n, k int
	a, b string
}

func genTestsG() []testG {
	rand.Seed(1530007)
	tests := make([]testG, 100)
	for i := range tests {
		n := rand.Intn(6) + 1
		k := rand.Intn(n + 1)
		ab := make([]byte, n)
		bb := make([]byte, n)
		for j := 0; j < n; j++ {
			ab[j] = byte('0' + rand.Intn(2))
			bb[j] = byte('0' + rand.Intn(2))
		}
		tests[i] = testG{n: n, k: k, a: string(ab), b: string(bb)}
	}
	return tests
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", oracle, "1530G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTestsG()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n%s\n%s\n", tc.n, tc.k, tc.a, tc.b)
	}

	cmdO := exec.Command(oracle)
	cmdO.Stdin = bytes.NewReader(input.Bytes())
	outO, err := cmdO.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
		os.Exit(1)
	}
	expected := strings.TrimSpace(string(outO))

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	outB, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, outB)
		os.Exit(1)
	}
	got := strings.TrimSpace(string(outB))
	if got != expected {
		fmt.Println("wrong answer")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
