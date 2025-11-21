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
	refSource2066C = "2066C.go"
	refBinary2066C = "ref2066C.bin"
	maxTests       = 160
	maxTotalN      = 200000
)

type testCase struct {
	a []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	cmd := exec.Command("go", "build", "-o", refBinary2066C, refSource2066C)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2066C), nil
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
		fmt.Fprintf(&sb, "%d\n", len(tc.a))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2066))
	var tests []testCase
	totalN := 0

	add := func(vals []int64) {
		tests = append(tests, testCase{a: vals})
		totalN += len(vals)
	}

	add([]int64{1})
	add([]int64{7, 9})
	add([]int64{1, 2, 3, 3, 2})
	add([]int64{1, 0, 0, 0})
	add([]int64{1000000000})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		n := rnd.Intn(min(remain, 3000)) + 1
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			if rnd.Intn(5) == 0 {
				vals[i] = int64(rnd.Intn(4))
			} else {
				vals[i] = rnd.Int63n(1_000_000_000) + 1
			}
		}
		add(vals)
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
