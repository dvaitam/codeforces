package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	refSource   = "2128B.go"
	totalNLimit = 180000
)

type testCase struct {
	p []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	// Validate our checker using the reference output.
	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s\n", err, refOut)
		os.Exit(1)
	}
	refSeqs, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s\n", err, refOut)
		os.Exit(1)
	}
	for i, seq := range refSeqs {
		if err := validateSequence(tests[i], seq); err != nil {
			fmt.Fprintf(os.Stderr, "internal error: reference output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	// Run candidate.
	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s\n", err, candOut)
		os.Exit(1)
	}
	candSeqs, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i, seq := range candSeqs {
		if err := validateSequence(tests[i], seq); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "permutation=%v\n", summarizePermutation(tests[i].p))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2128B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Join(dir, refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
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

func parseOutputs(out string, expected int) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]string, 0, expected)
	for sc.Scan() {
		res = append(res, sc.Text())
		if len(res) == expected {
			break
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}
	if len(res) < expected {
		return nil, fmt.Errorf("expected %d strings, got %d", expected, len(res))
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected starting at %q", sc.Text())
	}
	return res, nil
}

func validateSequence(tc testCase, seq string) error {
	n := len(tc.p)
	if len(seq) != n {
		return fmt.Errorf("sequence length %d does not match n=%d", len(seq), n)
	}
	l, r := 0, n-1
	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		ch := seq[i]
		if ch >= 'a' && ch <= 'z' {
			ch = ch - 'a' + 'A'
		}
		switch ch {
		case 'L':
			q = append(q, tc.p[l])
			l++
		case 'R':
			q = append(q, tc.p[r])
			r--
		default:
			return fmt.Errorf("invalid character %q at position %d", seq[i], i+1)
		}
	}
	if l != r+1 {
		return fmt.Errorf("did not consume all elements: l=%d r=%d", l, r)
	}
	if !isGood(q) {
		return fmt.Errorf("resulting array is bad: %v", summarizePermutation(q))
	}
	return nil
}

func isGood(arr []int) bool {
	n := len(arr)
	if n < 5 {
		return true
	}
	for i := 0; i+4 < n; i++ {
		inc, dec := true, true
		for j := 0; j < 4; j++ {
			if arr[i+j] >= arr[i+j+1] {
				inc = false
			}
			if arr[i+j] <= arr[i+j+1] {
				dec = false
			}
		}
		if inc || dec {
			return false
		}
	}
	return true
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.p)))
		sb.WriteByte('\n')
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	total := 0
	add := func(tc testCase) {
		if total+len(tc.p) > totalNLimit {
			return
		}
		tests = append(tests, tc)
		total += len(tc.p)
	}

	for _, tc := range deterministicTests() {
		add(tc)
	}

	// Large random stress.
	if total+100000 <= totalNLimit {
		add(randomPermutation(100000, rand.New(rand.NewSource(2128))))
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < totalNLimit {
		n := 5 + rng.Intn(2000)
		if total+n > totalNLimit {
			n = totalNLimit - total
		}
		add(randomPermutation(n, rng))
	}

	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{p: []int{1, 2, 3, 4, 5}},             // strictly increasing
		{p: []int{5, 4, 3, 2, 1}},             // strictly decreasing
		{p: []int{1, 2, 3, 4, 5, 6, 7}},       // sample-like
		{p: []int{1, 3, 6, 8, 9, 7, 5, 4, 2}}, // sample-like
		{p: []int{1, 2, 11, 3, 6, 4, 7, 8, 12, 5, 10, 9}},
		{p: []int{4, 1, 2, 5, 6, 3}},
		{p: []int{1, 2, 3, 5, 4}},
		{p: []int{5, 1, 8, 6, 2, 7, 9, 4, 3}},
		{p: []int{9, 2, 7, 5, 3, 8, 1, 6, 4, 10}},
	}
}

func randomPermutation(n int, rng *rand.Rand) testCase {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return testCase{p: arr}
}

func summarizePermutation(p []int) string {
	const limit = 15
	if len(p) <= limit {
		return fmt.Sprint(p)
	}
	head := fmt.Sprint(p[:limit])
	return head[:len(head)-1] + fmt.Sprintf(" ... total=%d]", len(p))
}
