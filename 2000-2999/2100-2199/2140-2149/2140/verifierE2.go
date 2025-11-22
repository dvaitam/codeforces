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
	refSource   = "2000-2999/2100-2199/2140-2149/2140/2140E2.go"
	randomTests = 120
	maxTotalExp = 1 << 20 // sum of 2^n across tests
	maxTotalM   = 800000  // keep below 1e6 limit
	maxN        = 20
	maxM        = 1_000_000
)

type testCase struct {
	n    int
	m    int
	k    int
	good []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/candidate")
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
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			tc := tests[i]
			fail("mismatch on test %d: expected %d, got %d (n=%d m=%d k=%d good=%v)", i+1, expect[i], got[i], tc.n, tc.m, tc.k, tc.good)
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2140E2-ref-*")
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

	totalExp := 0
	totalM := 0

	add := func(tc testCase) {
		exp := 1 << tc.n
		if tc.n < 0 || tc.n > maxN || exp <= 0 {
			return
		}
		if totalExp+exp > maxTotalExp {
			return
		}
		if totalM+tc.m > maxTotalM {
			return
		}
		tests = append(tests, tc)
		totalExp += exp
		totalM += tc.m
	}

	// Deterministic/sample-like cases
	add(testCase{n: 2, m: 3, k: 1, good: []int{1}})
	add(testCase{n: 7, m: 4, k: 3, good: []int{1, 4, 6}})
	add(testCase{n: 12, m: 3, k: 6, good: []int{1, 3, 5, 7, 9, 11}})
	add(testCase{n: 11, m: 12, k: 11, good: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}})
	add(testCase{n: 9, m: 6969, k: 2, good: []int{1, 9}})

	for len(tests) < randomTests && totalExp < maxTotalExp && totalM < maxTotalM {
		n := rng.Intn(maxN) + 1
		exp := 1 << n
		if totalExp+exp > maxTotalExp {
			continue
		}
		remainM := maxTotalM - totalM
		if remainM <= 0 {
			break
		}
		m := rng.Intn(minInt(maxM, remainM)) + 1

		k := rng.Intn(n) + 1
		good := make([]int, k)
		good[0] = 1
		for i := 1; i < k; i++ {
			good[i] = rng.Intn(n) + 1
		}
		// ensure sorted unique with 1 included
		good = uniqueSorted(good)
		k = len(good)

		add(testCase{n: n, m: m, k: k, good: good})
	}

	return tests
}

func uniqueSorted(arr []int) []int {
	m := make(map[int]struct{}, len(arr))
	m[1] = struct{}{}
	for _, v := range arr {
		if v < 1 {
			v = 1
		}
		if v > maxN {
			v = maxN
		}
		m[v] = struct{}{}
	}
	res := make([]int, 0, len(m))
	for v := range m {
		res = append(res, v)
	}
	for i := 0; i < len(res); i++ {
		for j := i + 1; j < len(res); j++ {
			if res[j] < res[i] {
				res[i], res[j] = res[j], res[i]
			}
		}
	}
	return res
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		sb.WriteString(strconv.Itoa(tc.k))
		sb.WriteByte('\n')
		for i, v := range tc.good {
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

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
