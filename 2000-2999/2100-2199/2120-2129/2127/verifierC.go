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
	refSource   = "2000-2999/2100-2199/2120-2129/2127/2127C.go"
	randomTests = 120
	maxTotalN   = 200000
	maxCaseN    = 60000
	maxValue    = 1_000_000_000
)

type testCase struct {
	n int
	k int
	a []int64
	b []int64
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
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			tc := tests[i]
			fail("mismatch on test %d: expected %d, got %d (n=%d, k=%d)", i+1, expect[i], got[i], tc.n, tc.k)
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2127C-ref-*")
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
		if tc.n <= 0 {
			return
		}
		if totalN+tc.n > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic coverage and sample-like scenarios.
	add(testCase{n: 2, k: 1, a: []int64{1, 7}, b: []int64{3, 5}})
	add(testCase{n: 3, k: 2, a: []int64{1, 5, 3}, b: []int64{6, 2, 4}})
	add(testCase{n: 3, k: 2, a: []int64{2, 2, 15}, b: []int64{15, 4, 12}})
	add(testCase{n: 4, k: 1, a: []int64{1, 6, 10, 10}, b: []int64{16, 3, 2, 10}})
	add(testCase{n: 2, k: 2, a: []int64{5, 5}, b: []int64{1, 9}})
	add(testCase{n: 5, k: 3, a: []int64{1, 2, 3, 4, 5}, b: []int64{5, 4, 3, 2, 1}})

	for len(tests) < randomTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		maxLen := maxCaseN
		if maxLen > remain {
			maxLen = remain
		}
		if maxLen < 2 {
			break
		}
		n := rng.Intn(maxLen-1) + 2 // at least 2
		k := rng.Intn(n) + 1

		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = rngInt64(rng, maxValue)
			b[i] = rngInt64(rng, maxValue)
		}

		// Occasionally create nearly equal arrays to test zero/low result.
		if rng.Intn(5) == 0 {
			val := rngInt64(rng, maxValue)
			for i := 0; i < n; i++ {
				a[i] = val
				if rng.Intn(3) == 0 {
					b[i] = val
				} else {
					b[i] = val + int64(rng.Intn(5))
				}
			}
		}

		add(testCase{n: n, k: k, a: a, b: b})
	}

	return tests
}

func rngInt64(rng *rand.Rand, limit int) int64 {
	return int64(rng.Intn(limit) + 1)
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.k))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
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

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
