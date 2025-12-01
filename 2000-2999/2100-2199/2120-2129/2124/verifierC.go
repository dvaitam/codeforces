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
	refSource   = "./2124C.go"
	targetTests = 120
	maxTotalN   = 200000
	maxValue    = 1_000_000_000
)

type testCase struct {
	b []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	// Sanity check reference answers.
	for i, x := range refAns {
		if !validX(tests[i].b, x) {
			fmt.Fprintf(os.Stderr, "reference produced invalid x on test %d\n", i+1)
			os.Exit(1)
		}
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i, x := range candAns {
		if !validX(tests[i].b, x) {
			fmt.Fprintf(os.Stderr, "test %d invalid x %d\n", i+1, x)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2124C-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func parseAnswers(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d: %v", tok, i+1, err)
		}
		if val <= 0 || val > maxValue {
			return nil, fmt.Errorf("x out of bounds at position %d: %d", i+1, val)
		}
		res[i] = val
	}
	return res, nil
}

func validX(b []int64, x int64) bool {
	if x <= 0 {
		return false
	}
	// options for first element.
	opts := optionsFor(b[0], x)
	if len(opts) == 0 {
		return false
	}
	for i := 1; i < len(b); i++ {
		curOpts := optionsFor(b[i], x)
		if len(curOpts) == 0 {
			return false
		}
		next := make([]int64, 0, 2)
		addIfMissing := func(v int64) {
			for _, existing := range next {
				if existing == v {
					return
				}
			}
			next = append(next, v)
		}
		for _, prev := range opts {
			for _, cand := range curOpts {
				if cand%prev == 0 {
					addIfMissing(cand)
				}
			}
		}
		if len(next) == 0 {
			return false
		}
		opts = next
	}
	return true
}

func optionsFor(val int64, x int64) []int64 {
	opts := []int64{val}
	if val%x == 0 {
		div := val / x
		if div != val {
			opts = append(opts, div)
		}
	}
	return opts
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", len(tc.b))
		for i, v := range tc.b {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		if totalN+len(tc.b) > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += len(tc.b)
	}

	// Sample-inspired tests.
	add(makeValidCase([]int64{1, 2}, 3, []int{0}))
	add(makeValidCase([]int64{2, 2, 2}, 2, []int{0, 1, 2}))
	add(makeValidCase([]int64{1, 2, 4, 8}, 4, []int{0, 1}))

	// Small n cases.
	add(randomCase(2, rng))
	add(randomCase(3, rng))

	// Random tests until limits.
	for len(tests) < targetTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain < 2 {
			break
		}
		n := rng.Intn(min(2000, remain-1)) + 2
		add(randomCase(n, rng))
	}

	if len(tests) == 0 {
		add(randomCase(2, rng))
	}
	return tests
}

func randomCase(n int, rng *rand.Rand) testCase {
	// Build a beautiful array a then multiply subset by x.
	a := make([]int64, n)
	a[0] = rng.Int63n(1_000) + 1
	for i := 1; i < n; i++ {
		mult := int64(rng.Intn(3) + 1)
		val := a[i-1] * mult
		if val > maxValue/10 {
			val = a[i-1]
		}
		a[i] = val
	}
	maxA := a[0]
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
	}
	maxX := maxValue / maxA
	if maxX < 1 {
		maxX = 1
	}
	if maxX > 1000 {
		maxX = 1000
	}
	x := int64(rng.Intn(int(maxX)) + 1)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = a[i] * x
		} else {
			b[i] = a[i]
		}
	}
	return testCase{b: b}
}

func makeValidCase(a []int64, x int64, subset []int) testCase {
	b := make([]int64, len(a))
	isSub := make(map[int]struct{})
	for _, idx := range subset {
		isSub[idx] = struct{}{}
	}
	for i, v := range a {
		if _, ok := isSub[i]; ok {
			b[i] = v * x
		} else {
			b[i] = v
		}
	}
	return testCase{b: b}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
