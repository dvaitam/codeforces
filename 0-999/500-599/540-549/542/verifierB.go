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
	"time"
)

const (
	refSource542B = "542B.go"
	refBinary542B = "ref542B.bin"
	maxTests      = 120
	maxTotalN     = 200000
)

type testCase struct {
	n        int
	r        int64
	segments [][2]int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary542B, refSource542B)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary542B), nil
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
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
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
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.r)
		for _, iv := range tc.segments {
			fmt.Fprintf(&sb, "%d %d\n", iv[0], iv[1])
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	add(testCase{
		n: 3,
		r: 3,
		segments: [][2]int64{
			{-3, 0},
			{1, 3},
			{-1, 2},
		},
	})
	add(testCase{
		n: 4,
		r: 5,
		segments: [][2]int64{
			{-1, 1},
			{2, 4},
			{5, 9},
			{6, 8},
		},
	})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain <= 0 {
			break
		}
		n := rnd.Intn(2000) + 1
		if n > remain {
			n = remain
		}
		r := rnd.Int63n(1_000_000_000) + 1
		segs := make([][2]int64, n)
		for i := 0; i < n; i++ {
			h := rnd.Int63n(2_000_000_001) - 1_000_000_000
			length := rnd.Int63n(1_000_000_000) + 1
			t := h + length
			segs[i] = [2]int64{h, t}
		}
		add(testCase{n: n, r: r, segments: segs})
	}
	return tests
}
