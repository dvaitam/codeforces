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

type testCase struct {
	id    string
	input string
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}

	base := currentDir()
	refBin, err := buildReference(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "wrong answer on %s\nInput:\n%sExpected: %s\nGot: %s\n", tc.id, tc.input, exp, got)
			os.Exit(1)
		}
		if (i+1)%10 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d/%d tests...\n", i+1, len(tests))
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot determine current file path")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref2147G.bin")
	cmd := exec.Command("go", "build", "-o", out, "2147G.go")
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	var tests []testCase
	deterministic := [][]int64{
		{5, 1, 1},
		{5, 2, 1},
		{23, 1, 1},
		{10, 10, 2},
		{999, 1, 1},
		{1, 1, 2},
	}
	tests = append(tests, makeCase("samples", deterministic))

	// pairwise coprime, repeated primes, etc.
	tests = append(tests, makeCase("prime-power", [][]int64{
		{2, 4, 8},
		{3, 9, 27},
		{7, 7, 7},
		{11, 11, 13},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("rand-%02d", i+1), rng.Intn(5)+1))
	}
	// stress near limit
	for i := 0; i < 5; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("stress-%02d", i+1), 100))
	}
	return tests
}

func makeCase(id string, triples [][]int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(triples))
	for _, tri := range triples {
		fmt.Fprintf(&sb, "%d %d %d\n", tri[0], tri[1], tri[2])
	}
	return testCase{id: id, input: sb.String()}
}

func randomCase(rng *rand.Rand, id string, t int) testCase {
	var triples [][]int64
	for i := 0; i < t; i++ {
		x := int64(rng.Intn(1_000_000) + 1)
		y := int64(rng.Intn(1_000_000) + 1)
		z := int64(rng.Intn(1_000_000) + 1)
		triples = append(triples, []int64{x, y, z})
	}
	return makeCase(id, triples)
}
