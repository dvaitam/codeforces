package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource2115F2 = "./2115F2.go"

type op struct {
	a int
	b int
	c int
}

type testCase struct {
	n   int
	q   int
	ops []op
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s",
				idx+1, err, input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\ninput:\n%soutput:\n%s",
				idx+1, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s",
				idx+1, err, input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\ninput:\n%soutput:\n%s",
				idx+1, err, input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.q; i++ {
			if refAns[i] != candAns[i] {
				fmt.Fprintf(os.Stderr, "Mismatch on test %d, query %d: expected %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, i+1, refAns[i], candAns[i], input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2115F2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2115F2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2115F2)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.Write(errBuf.Bytes())
	}
	return out.String(), err
}

func parseOutput(out string, q int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != q {
		return nil, fmt.Errorf("expected %d outputs, got %d", q, len(fields))
	}
	res := make([]int, q)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
	for _, op := range tc.ops {
		fmt.Fprintf(&sb, "%d %d %d\n", op.a, op.b, op.c)
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 30)

	// Small deterministic scenarios.
	tests = append(tests,
		testCase{
			n: 1, q: 5,
			ops: []op{
				{a: 1, b: 1, c: 1},
				{a: 1, b: 1, c: 1},
				{a: 2, b: 1, c: 1},
				{a: 3, b: 1, c: 1},
				{a: 3, b: 1, c: 1},
			},
		},
		testCase{
			n: 3, q: 6,
			ops: []op{
				{a: 1, b: 2, c: 1},
				{a: 2, b: 3, c: 2},
				{a: 1, b: 3, c: 3},
				{a: 3, b: 2, c: 2},
				{a: 2, b: 1, c: 3},
				{a: 1, b: 3, c: 1},
			},
		},
		testCase{
			n: 5, q: 8,
			ops: []op{
				{a: 1, b: 5, c: 3},
				{a: 2, b: 4, c: 4},
				{a: 1, b: 2, c: 5},
				{a: 3, b: 1, c: 1},
				{a: 3, b: 8, c: 2},
				{a: 2, b: 5, c: 3},
				{a: 1, b: 1, c: 4},
				{a: 3, b: 4, c: 5},
			},
		},
	)

	rng := rand.New(rand.NewSource(2115_0202))
	for len(tests) < 25 {
		n := rng.Intn(40) + 1  // 1..40
		q := rng.Intn(80) + 20 // 20..99
		ops := make([]op, q)
		for i := 0; i < q; i++ {
			a := rng.Intn(3) + 1
			var b int
			if a == 3 {
				b = rng.Intn(q) + 1
			} else {
				b = rng.Intn(n) + 1
			}
			c := rng.Intn(n) + 1
			ops[i] = op{a: a, b: b, c: c}
		}
		tests = append(tests, testCase{n: n, q: q, ops: ops})
	}

	return tests
}
