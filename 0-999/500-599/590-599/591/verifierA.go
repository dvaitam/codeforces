package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	l int
	p int
	q int
}

func generateTests() []testCase {
	tests := []testCase{
		{1, 1, 1},
		{1000, 500, 500},
		{1000, 1, 500},
		{500, 500, 1},
		{999, 123, 456},
	}
	rng := rand.New(rand.NewSource(1))
	for len(tests) < 100 {
		l := rng.Intn(1000) + 1
		p := rng.Intn(500) + 1
		q := rng.Intn(500) + 1
		tests = append(tests, testCase{l, p, q})
	}
	return tests
}

func solve(tc testCase) float64 {
	return float64(tc.p) * float64(tc.l) / float64(tc.p+tc.q)
}

func run(bin string, input string) (string, error) {
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

func equal(a, b float64) bool {
	diff := math.Abs(a - b)
	return diff <= 1e-4*math.Max(1, math.Abs(b))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d\n%d\n%d\n", tc.l, tc.p, tc.q)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Sscanf(gotStr, "%f", &got); err != nil {
			fmt.Fprintf(os.Stderr, "cannot parse output on test %d: %v\noutput: %s\n", i+1, err, gotStr)
			os.Exit(1)
		}
		exp := solve(tc)
		if !equal(got, exp) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected: %.6f\ngot: %.6f\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
