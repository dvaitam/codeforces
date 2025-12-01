package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	refSource    = "./2118D2.go"
	randomTests  = 80
	maxTotalSize = 180000 // sum of n+q to keep runtime comfortable (limit 2e5)
	maxN         = 200000
	maxQ         = 200000
)

type testCase struct {
	n int
	k int64
	p []int64
	d []int64
	q int
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/candidate")
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
	expect, err := parseOutputs(expectRaw, tests)
	if err != nil {
		fail("could not parse reference output: %v\n%s", err, expectRaw)
	}

	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}
	got, err := parseOutputs(gotRaw, tests)
	if err != nil {
		fail("could not parse candidate output: %v\n%s", err, gotRaw)
	}

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d answers, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at answer %d: expected %s, got %s", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2118D2-ref-*")
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
	tests := make([]testCase, 0, randomTests+4)

	totalSize := 0
	add := func(tc testCase) {
		if tc.n <= 0 || tc.q <= 0 {
			return
		}
		addSize := tc.n + tc.q
		if totalSize+addSize > maxTotalSize {
			return
		}
		tests = append(tests, tc)
		totalSize += addSize
	}

	// Deterministic small/edge cases.
	add(makeCase(2, 2, []int64{1, 4}, []int64{1, 0}, []int64{1, 2, 3}))
	add(makeCase(1, 5, []int64{3}, []int64{0}, []int64{1, 2, 10}))
	add(makeCase(3, 1, []int64{2, 5, 7}, []int64{0, 0, 0}, []int64{1, 6, 8}))
	add(makeCase(4, 3, []int64{3, 4, 5, 6}, []int64{0, 1, 2, 0}, []int64{3, 7, 9}))

	for len(tests) < randomTests && totalSize < maxTotalSize {
		remaining := maxTotalSize - totalSize
		// Small bias toward smaller n and q to diversify.
		n := rng.Intn(minInt(maxN, max(1, remaining/2))) + 1
		if n > 20 && rng.Intn(4) == 0 {
			n = rng.Intn(30) + 1
		}
		maxQForRemain := remaining - n
		if maxQForRemain < 1 {
			break
		}
		q := rng.Intn(minInt(maxQ, maxQForRemain)) + 1

		k := rng.Int63n(1_000_000_000_000_000) + 1 // up to 1e15
		p := make([]int64, n)
		curr := rng.Int63n(k) + 1
		for i := 0; i < n; i++ {
			// ensure increasing positions within 1..1e15
			step := rng.Int63n(1_000_000) + 1
			curr += step
			if curr > 1_000_000_000_000_000 {
				curr = 1_000_000_000_000_000
			}
			p[i] = curr
		}
		sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })
		d := make([]int64, n)
		for i := 0; i < n; i++ {
			if k == 1 {
				d[i] = 0
			} else {
				d[i] = rng.Int63n(k)
			}
		}

		a := make([]int64, q)
		for i := 0; i < q; i++ {
			// choose around positions and random far points
			if rng.Intn(3) == 0 {
				a[i] = p[rng.Intn(n)]
			} else {
				a[i] = rng.Int63n(1_000_000_000_000_000) + 1
			}
		}

		add(testCase{n: n, k: k, p: p, d: d, q: q, a: a})
	}

	return tests
}

func makeCase(n int, k int64, p []int64, d []int64, a []int64) testCase {
	if len(p) != n || len(d) != n {
		return testCase{}
	}
	if !sort.SliceIsSorted(p, func(i, j int) bool { return p[i] < p[j] }) {
		sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })
	}
	return testCase{n: n, k: k, p: p, d: d, q: len(a), a: a}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.d {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(tc.q))
		sb.WriteByte('\n')
		for i, v := range tc.a {
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

func parseOutputs(out string, tests []testCase) ([]string, error) {
	lines := strings.Fields(out)
	expectCount := 0
	for _, tc := range tests {
		expectCount += tc.q
	}
	if len(lines) != expectCount {
		return nil, fmt.Errorf("expected %d answers, got %d", expectCount, len(lines))
	}
	res := make([]string, expectCount)
	for i, s := range lines {
		up := strings.ToUpper(s)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid answer %q at position %d", s, i+1)
		}
		res[i] = up
	}
	return res, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
