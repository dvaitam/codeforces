package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const refSource = "2110B.go"

type testCase struct {
	s string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\n", i+1, refAns[i], candAns[i])
			fmt.Fprintf(os.Stderr, "s=%s\n", tests[i].s)
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
	tmp, err := os.CreateTemp("", "ref_2110B_*.bin")
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
	add := func(s string) {
		tests = append(tests, testCase{s: s})
	}

	// Fixed examples and edge cases
	add("(())")
	add("(())()()()")
	add("()")
	add("()()")
	add("((()))")
	add("((()))()")
	add("()()()")
	add("(()())")

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomBalanced := func(length int) string {
		pairs := length / 2
		openRem, closeRem := pairs, pairs
		openCount := 0
		var sb strings.Builder
		for i := 0; i < 2*pairs; i++ {
			// must keep sequence valid
			if openRem == 0 {
				sb.WriteByte(')')
				closeRem--
				openCount--
				continue
			}
			if closeRem == openCount {
				sb.WriteByte('(')
				openRem--
				openCount++
				continue
			}
			if openCount == 0 {
				sb.WriteByte('(')
				openRem--
				openCount++
				continue
			}
			if rng.Intn(2) == 0 {
				sb.WriteByte('(')
				openRem--
				openCount++
			} else {
				sb.WriteByte(')')
				closeRem--
				openCount--
			}
		}
		return sb.String()
	}

	// Random cases with varying lengths
	for len(tests) < 150 {
		length := rng.Intn(40) + 2
		if length%2 == 1 {
			length++
		}
		add(randomBalanced(length))
	}

	// One large stress case near the limit
	largeLen := 200000
	add(randomBalanced(largeLen))

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	res := make([]string, expected)
	for i, s := range lines {
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
