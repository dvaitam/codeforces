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
	refSource   = "2000-2999/2100-2199/2130-2139/2132/2132G.go"
	randomTests = 100
	maxTotalNM  = 800000 // below 1e6 limit to stay fast
	maxN        = 1000000
	maxM        = 1000000
)

type testCase struct {
	n    int
	m    int
	grid []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
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
			fail("mismatch on test %d: expected %d, got %d (n=%d m=%d)", i+1, expect[i], got[i], tc.n, tc.m)
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2132G-ref-*")
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

	totalNM := 0
	add := func(tc testCase) {
		if tc.n <= 0 || tc.m <= 0 {
			return
		}
		area := tc.n * tc.m
		if area == 0 || totalNM+area > maxTotalNM {
			return
		}
		tests = append(tests, tc)
		totalNM += area
	}

	// Deterministic cases including sample-like shapes and edge cases.
	add(makeGrid([]string{"hey", "hey"}))
	add(makeGrid([]string{"abc", "def", "ghi"}))
	add(makeGrid([]string{"af", "fa", "te"}))
	add(makeGrid([]string{"x"}))
	add(makeGrid([]string{"uoe", "vbe", "mbu"}))
	add(makeGrid([]string{"hyh", "kop"}))
	add(makeGrid([]string{"a", "b"}))
	add(makeGrid([]string{"aaaa", "aaaa"}))

	for len(tests) < randomTests && totalNM < maxTotalNM {
		remain := maxTotalNM - totalNM
		// bias to smaller grids but allow occasional larger
		maxSide := 1000
		if maxSide*maxSide > remain {
			maxSide = int(sqrtInt(remain)) + 1
		}
		if maxSide < 1 {
			break
		}
		n := rng.Intn(maxSide) + 1
		mLimit := remain / n
		if mLimit < 1 {
			break
		}
		mCap := minInt(maxSide, mLimit)
		m := rng.Intn(mCap) + 1

		grid := make([]string, n)
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				row[j] = byte('a' + rng.Intn(26))
			}
			grid[i] = string(row)
		}
		add(testCase{n: n, m: m, grid: grid})
	}

	return tests
}

func makeGrid(rows []string) testCase {
	n := len(rows)
	m := len(rows[0])
	for _, s := range rows {
		if len(s) != m {
			return testCase{}
		}
	}
	return testCase{n: n, m: m, grid: rows}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
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

func sqrtInt(x int) int {
	r := int(1)
	for (r+1)*(r+1) <= x {
		r++
	}
	return r
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
