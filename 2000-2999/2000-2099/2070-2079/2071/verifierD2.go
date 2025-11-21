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

const (
	refSource2071D2 = "2071D2.go"
	refBinary2071D2 = "ref2071D2.bin"
	maxTests        = 120
	maxTotalN       = 200000
)

type testCase struct {
	n int
	l int64
	r int64
	a []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on case %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2071D2, refSource2071D2)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2071D2), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.l, tc.r)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte(byte('0' + v))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2071))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	add(testCase{n: 1, l: 1, r: 1, a: []int{1}})
	add(testCase{n: 2, l: 1, r: 3, a: []int{0, 1}})
	add(testCase{n: 5, l: 3, r: 10, a: []int{1, 0, 1, 1, 0}})
	add(testCase{n: 6, l: 1, r: 1_000_000_000_000_000_000, a: []int{1, 0, 1, 0, 1, 0}})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		n := rnd.Intn(min(remain, 40000)) + 1
		l := rnd.Int63n(1_000_000_000_000_000_000) + 1
		r := l + rnd.Int63n(1_000_000_000_000_000_000-l+1)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if rnd.Intn(2) == 0 {
				a[i] = 0
			} else {
				a[i] = 1
			}
		}
		add(testCase{n: n, l: l, r: r, a: a})
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
