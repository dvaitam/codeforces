package main

import (
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

const refSource = "2144A.go"

type testCase struct {
	n int
	a []int
}

type pair struct {
	l int
	r int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refPairs, err := parsePairs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	refChecked := true
	for i, tc := range tests {
		exists := existsSolution(tc)
		if err := checkPair(tc, refPairs[i], exists); err != nil {
			refChecked = false
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	if !refChecked {
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candPairs, err := parsePairs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		exists := existsSolution(tc)
		if err := checkPair(tc, candPairs[i], exists); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "n=%d a=%v\n", tc.n, tc.a)
			fmt.Fprintf(os.Stderr, "expected existence=%v, candidate pair=(%d,%d)\n", exists, candPairs[i].l, candPairs[i].r)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	refPath, err := referencePath()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "ref_2144A_*.bin")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate verifier path")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, refSource), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(arr []int) {
		tests = append(tests, testCase{n: len(arr), a: arr})
	}

	// Fixed cases including sample-like scenarios
	add([]int{1, 2, 3, 4, 5, 6})
	add([]int{1, 3, 3, 7})
	add([]int{2, 1, 0})
	add([]int{7, 2, 6, 2, 4})
	add([]int{0, 0, 0})
	add([]int{3, 3, 3, 3})
	add([]int{1, 2, 3})
	add([]int{5, 5, 5, 5, 5})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		n := rng.Intn(38) + 3 // 3..40
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(41)
		}
		add(arr)
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parsePairs(out string, expected int) ([]pair, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected*2 {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected*2, len(tokens))
	}
	res := make([]pair, expected)
	for i := 0; i < expected; i++ {
		l, err := strconv.Atoi(tokens[2*i])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[2*i])
		}
		r, err := strconv.Atoi(tokens[2*i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[2*i+1])
		}
		res[i] = pair{l: l, r: r}
	}
	return res, nil
}

func existsSolution(tc testCase) bool {
	n := tc.n
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + tc.a[i-1]
	}
	for l := 1; l <= n-2; l++ {
		s1 := pref[l] % 3
		for r := l + 1; r <= n-1; r++ {
			s2 := (pref[r] - pref[l]) % 3
			if s2 < 0 {
				s2 += 3
			}
			s3 := (pref[n] - pref[r]) % 3
			if s3 < 0 {
				s3 += 3
			}
			if (s1 == s2 && s2 == s3) || (s1 != s2 && s1 != s3 && s2 != s3) {
				return true
			}
		}
	}
	return false
}

func checkPair(tc testCase, pr pair, exists bool) error {
	n := tc.n
	if pr.l == 0 && pr.r == 0 {
		if exists {
			return fmt.Errorf("solution exists but got 0 0")
		}
		return nil
	}
	if pr.l < 1 || pr.r < 1 || pr.l >= pr.r || pr.r >= n {
		return fmt.Errorf("invalid indices (%d,%d)", pr.l, pr.r)
	}
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + tc.a[i-1]
	}
	s1 := pref[pr.l] % 3
	s2 := (pref[pr.r] - pref[pr.l]) % 3
	if s2 < 0 {
		s2 += 3
	}
	s3 := (pref[n] - pref[pr.r]) % 3
	if s3 < 0 {
		s3 += 3
	}
	ok := (s1 == s2 && s2 == s3) || (s1 != s2 && s1 != s3 && s2 != s3)
	if !ok {
		return fmt.Errorf("pair (%d,%d) does not satisfy condition", pr.l, pr.r)
	}
	if !exists {
		return fmt.Errorf("reported solution but none exists")
	}
	return nil
}
