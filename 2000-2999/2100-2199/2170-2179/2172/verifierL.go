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

type testInput struct {
	name string
	n    int
	m    int
	k    int
	s    string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2172L-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleL")
	cmd := exec.Command("go", "build", "-o", bin, "2172L.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testInput) string {
	return fmt.Sprintf("%d %d %d\n%s\n", tc.n, tc.m, tc.k, tc.s)
}

func deterministicTests() []testInput {
	return []testInput{
		{name: "sample1", n: 5, m: 4, k: 3, s: "RRBRR"},
		{name: "sample2", n: 10, m: 3, k: 3, s: "RRRRBBRRRB"},
		{name: "sample3", n: 7, m: 4, k: 7, s: "RRBRBBR"},
		{name: "single_char", n: 1, m: 0, k: 1, s: "R"},
		{name: "k_equals_n", n: 6, m: 5, k: 6, s: "RBBRBR"},
		{name: "alternating", n: 8, m: 8, k: 2, s: "RBRBRBRB"},
		{name: "all_same_large_m", n: 12, m: 10, k: 5, s: strings.Repeat("R", 12)},
	}
}

func randomString(n int, rng *rand.Rand) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bytes[i] = 'R'
		} else {
			bytes[i] = 'B'
		}
	}
	return string(bytes)
}

func randomTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testInput, 0, 30)
	for i := 0; i < 20; i++ {
		n := rng.Intn(50) + 1
		m := rng.Intn(50)
		k := rng.Intn(n) + 1
		tests = append(tests, testInput{
			name: fmt.Sprintf("random_small_%d", i+1),
			n:    n,
			m:    m,
			k:    k,
			s:    randomString(n, rng),
		})
	}
	for i := 0; i < 10; i++ {
		n := rng.Intn(3000) + 1
		m := rng.Intn(3000)
		k := rng.Intn(n) + 1
		tests = append(tests, testInput{
			name: fmt.Sprintf("random_large_%d", i+1),
			n:    n,
			m:    m,
			k:    k,
			s:    randomString(n, rng),
		})
	}
	return tests
}

func stressTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 123))
	return []testInput{
		{
			name: "stress_all_R",
			n:    3000,
			m:    3000,
			k:    1,
			s:    strings.Repeat("R", 3000),
		},
		{
			name: "stress_all_B",
			n:    3000,
			m:    3000,
			k:    3000,
			s:    strings.Repeat("B", 3000),
		},
		{
			name: "stress_random",
			n:    3000,
			m:    3000,
			k:    1500,
			s:    randomString(3000, rng),
		},
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	tests = append(tests, stressTests()...)

	for idx, tc := range tests {
		input := buildInput(tc)
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		if expected != actual {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %s got %s\ninput:\n%s", idx+1, tc.name, expected, actual, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
