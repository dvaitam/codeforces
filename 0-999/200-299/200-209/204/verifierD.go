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
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "204D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

type testCase struct {
	n int
	k int
	s string
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(4))
	tests := make([]testCase, 0, 100)
	fixed := []testCase{{n: 1, k: 1, s: "B"}, {n: 3, k: 1, s: "BWB"}}
	tests = append(tests, fixed...)
	letters := []byte("BWX")
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		k := rng.Intn(n) + 1
		sb := make([]byte, n)
		for i := 0; i < n; i++ {
			sb[i] = letters[rng.Intn(3)]
		}
		tests = append(tests, testCase{n: n, k: k, s: string(sb)})
	}
	return tests
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
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n%s\n", t.n, t.k, t.s)
		expOut, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expOut {
			fmt.Printf("test %d failed. Expected %q got %q\n", i+1, expOut, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
