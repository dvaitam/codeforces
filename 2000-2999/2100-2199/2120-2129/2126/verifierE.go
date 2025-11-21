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

const refSource = "2126E.go"

type testCase struct {
	n int
	p []int
	s []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\n", i+1, refAns[i], candAns[i])
			fmt.Fprintf(os.Stderr, "n=%d\np=%v\ns=%v\n", tc.n, tc.p, tc.s)
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
	tmp, err := os.CreateTemp("", "ref_2126E_*.bin")
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
	add := func(tc testCase) {
		tests = append(tests, tc)
	}

	// Small deterministic cases
	add(makeFromArray([]int{72, 24, 3, 6, 12, 144}))
	add(makeFromArray([]int{125, 125, 125, 25, 75}))
	add(makeFromArray([]int{5}))
	add(makeFromArray([]int{2, 2}))

	// Construct invalid by tweaking valid cases
	invalid := makeFromArray([]int{10, 5, 5})
	invalid.p[1]++ // break prefix gcd chain
	add(invalid)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	for len(tests) < 120 && totalN < 100000 {
		n := rng.Intn(2000) + 1
		if totalN+n > 100000 {
			n = 100000 - totalN
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(1_000_000_000) + 1
		}
		tc := makeFromArray(a)

		// Occasionally perturb to make inconsistent inputs.
		if rng.Intn(3) == 0 {
			pos := rng.Intn(n)
			tc.p[pos] += int(rng.Int31n(3) + 1)
		}

		add(tc)
		totalN += n
	}

	return tests
}

func makeFromArray(a []int) testCase {
	n := len(a)
	p := make([]int, n)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			p[i] = a[i]
		} else {
			p[i] = gcd(p[i-1], a[i])
		}
	}
	for i := n - 1; i >= 0; i-- {
		if i == n-1 {
			s[i] = a[i]
		} else {
			s[i] = gcd(s[i+1], a[i])
		}
	}
	return testCase{n: n, p: p, s: s}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.s {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	res := make([]string, expected)
	for i, s := range fields {
		res[i] = normalizeAnswer(s)
	}
	return res, nil
}

func normalizeAnswer(s string) string {
	if strings.EqualFold(s, "YES") {
		return "YES"
	}
	return "NO"
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
