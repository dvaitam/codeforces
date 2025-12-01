package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource   = "./2115C.go"
	randomTests = 100
	maxTotalN   = 100
	maxN        = 20
	maxM        = 4000
	maxH        = 400
	tolerance   = 1e-6
)

type testCase struct {
	n int
	m int
	p int
	h []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	expect, err := parseOutputs(expectRaw, len(tests))
	if err != nil {
		fail("could not parse reference output: %v\n%s", err, expectRaw)
	}

	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}
	got, err := parseOutputs(gotRaw, len(tests))
	if err != nil {
		fail("could not parse candidate output: %v\n%s", err, gotRaw)
	}

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d lines, got %d", len(expect), len(got))
	}
	for i := range expect {
		if !closeEnough(expect[i], got[i]) {
			tc := tests[i]
			fail("mismatch on test %d: expected %.9f, got %.9f (n=%d m=%d p=%d h=%v)", i+1, expect[i], got[i], tc.n, tc.m, tc.p, tc.h)
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2115C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+6)

	totalN := 0
	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic coverage (sample-like and edges).
	add(testCase{n: 2, m: 2, p: 10, h: []int{2, 2}})
	add(testCase{n: 5, m: 5, p: 20, h: []int{2, 2, 2, 2, 2}})
	add(testCase{n: 6, m: 20, p: 50, h: []int{1, 1, 4, 5, 1, 4}})
	add(testCase{n: 1, m: 1, p: 0, h: []int{1}})
	add(testCase{n: 1, m: 5, p: 100, h: []int{400}})
	add(testCase{n: 3, m: 10, p: 0, h: []int{3, 3, 3}})
	add(testCase{n: 3, m: 10, p: 100, h: []int{3, 3, 3}})

	for len(tests) < randomTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		n := rng.Intn(min(maxN, remain)) + 1
		m := rng.Intn(maxM) + 1
		p := rng.Intn(101)

		h := make([]int, n)
		for i := 0; i < n; i++ {
			h[i] = rng.Intn(maxH) + 1
		}
		add(testCase{n: n, m: m, p: p, h: h})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.p))
		for i, v := range tc.h {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutputs(out string, t int) ([]float64, error) {
	lines := strings.Fields(out)
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(lines))
	}
	res := make([]float64, t)
	for i, tok := range lines {
		v, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func closeEnough(a, b float64) bool {
	diff := math.Abs(a - b)
	den := math.Max(1.0, math.Abs(b))
	return diff/den <= tolerance
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
