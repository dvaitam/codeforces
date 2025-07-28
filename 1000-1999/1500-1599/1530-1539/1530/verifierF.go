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

type testF struct {
	n int
	a [][]int
}

func genTestsF() []testF {
	rand.Seed(1530006)
	tests := make([]testF, 100)
	for i := range tests {
		n := rand.Intn(4) + 1 // small for speed
		mat := make([][]int, n)
		for j := 0; j < n; j++ {
			mat[j] = make([]int, n)
			for k := 0; k < n; k++ {
				mat[j][k] = rand.Intn(10001)
			}
		}
		tests[i] = testF{n: n, a: mat}
	}
	return tests
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1530F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTestsF()
	for idx, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintln(&input, tc.n)
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.n; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprint(&input, tc.a[i][j])
			}
			input.WriteByte('\n')
		}

		cmdO := exec.Command(oracle)
		cmdO.Stdin = bytes.NewReader(input.Bytes())
		outO, err := cmdO.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(string(outO))

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		outB, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s\n", idx+1, err, outB)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(outB))
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
