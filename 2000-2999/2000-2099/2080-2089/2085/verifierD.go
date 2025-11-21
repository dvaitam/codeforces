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
	refSourceD = "2085D.go"
	refBinaryD = "refD.bin"
	maxSumN    = 200000
	maxTests   = 200
)

type testCaseD struct {
	n int
	k int
	d []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	tcaseCount := len(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	refVals, err := parseOutput(refOut, tcaseCount)
	if err != nil {
		fmt.Printf("reference output parse error: %v\n", err)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	candVals, err := parseOutput(candOut, tcaseCount)
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := 0; i < tcaseCount; i++ {
		if refVals[i] != candVals[i] {
			fmt.Printf("Mismatch on test #%d: expected %d, got %d\n", i+1, refVals[i], candVals[i])
			fmt.Println("Input:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", tcaseCount)
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryD, refSourceD)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryD), nil
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
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d not an integer: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func formatInput(tests []testCaseD) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for i, val := range tc.d {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(val, 10))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func generateTests() []testCaseD {
	rnd := rand.New(rand.NewSource(2085))
	var tests []testCaseD
	totalN := 0

	prepend := []testCaseD{
		{n: 2, k: 1, d: []int64{1, 2}},
		{n: 3, k: 1, d: []int64{5, 1, 4}},
		{n: 5, k: 2, d: []int64{1, 2, 3, 4, 5}},
		{n: 6, k: 1, d: []int64{1, 1000000000, 1, 1000000000, 1, 1000000000}},
	}

	for _, tc := range prepend {
		tests = append(tests, tc)
		totalN += tc.n
	}

	for len(tests) < maxTests && totalN < maxSumN {
		maxN := maxSumN - totalN
		if maxN < 2 {
			break
		}
		n := rnd.Intn(min(5000, maxN-1)) + 2
		k := rnd.Intn(n-1) + 1
		d := make([]int64, n)
		for i := 0; i < n; i++ {
			if rnd.Float64() < 0.2 {
				d[i] = 1000000000
			} else {
				d[i] = int64(rnd.Intn(1000000000) + 1)
			}
		}
		tests = append(tests, testCaseD{n: n, k: k, d: d})
		totalN += n
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
