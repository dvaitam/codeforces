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

type testCase struct {
	input string
	n     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		refValid, refSeq, err := parseOutput(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\nOutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		candValid, candSeq, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\nOutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}
		if refValid == 0 {
			if candValid != 0 {
				fmt.Fprintf(os.Stderr, "candidate claims possible while reference says no on test %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, tc.input, candOut)
				os.Exit(1)
			}
			continue
		}
		if candValid == 0 {
			fmt.Fprintf(os.Stderr, "candidate claims impossible while reference provides a solution on test %d\nInput:\n%s\n", idx+1, tc.input)
			os.Exit(1)
		}
		if err := validateSequence(tc.input, candSeq); err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid sequence on test %d: %v\nInput:\n%s\nCandidate output:\n%s\n", idx+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1053E.go",
		filepath.Join("1000-1999", "1000-1099", "1050-1059", "1053", "1053E.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1053E.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1053E_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
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
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, n int) (int, []int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	first := strings.TrimSpace(strings.ToLower(lines[0]))
	switch first {
	case "no":
		return 0, nil, nil
	case "yes":
		// continue
	default:
		return 0, nil, fmt.Errorf("expected yes/no, got %q", lines[0])
	}
	if len(lines) < 2 {
		return 0, nil, fmt.Errorf("missing sequence line")
	}
	parts := strings.Fields(lines[1])
	if len(parts) != 2*n-1 {
		return 0, nil, fmt.Errorf("expected %d numbers, got %d", 2*n-1, len(parts))
	}
	seq := make([]int, len(parts))
	for i, tok := range parts {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid integer %q", tok)
		}
		seq[i] = val
	}
	return 1, seq, nil
}

func validateSequence(input string, seq []int) error {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return fmt.Errorf("failed to read n: %v", err)
	}
	targetLen := 2*n - 1
	if len(seq) != targetLen {
		return fmt.Errorf("sequence length mismatch: expected %d, got %d", targetLen, len(seq))
	}
	given := make([]int, targetLen)
	for i := 0; i < targetLen; i++ {
		if _, err := fmt.Fscan(reader, &given[i]); err != nil {
			return fmt.Errorf("failed to read given sequence: %v", err)
		}
	}
	if seq[0] != seq[len(seq)-1] {
		return fmt.Errorf("first and last elements differ")
	}
	for i := 0; i < len(seq)-1; i++ {
		if seq[i] == 0 || seq[i+1] == 0 || seq[i] == seq[i+1] {
			return fmt.Errorf("invalid edge between positions %d and %d", i+1, i+2)
		}
	}
	seen := make([]bool, n+1)
	for _, v := range seq {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range", v)
		}
		seen[v] = true
	}
	for i := 1; i <= n; i++ {
		if !seen[i] {
			return fmt.Errorf("vertex %d missing", i)
		}
	}
	for i := 0; i < targetLen; i++ {
		if given[i] != 0 && given[i] != seq[i] {
			return fmt.Errorf("position %d must be %d, got %d", i+1, given[i], seq[i])
		}
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		newTestCase(1, []int{0}),
		newTestCase(2, []int{1, 0, 1}),
		newTestCase(3, []int{0, 0, 0, 0, 0}),
		newTestCase(3, []int{1, 2, 1, 3, 1}),
		newTestCase(4, []int{0, 2, 0, 3, 0, 4, 0}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		tests = append(tests, randomValidTest(rng, rng.Intn(8)+1))
	}
	for len(tests) < 80 {
		tests = append(tests, randomValidTest(rng, rng.Intn(100)+5))
	}
	tests = append(tests, randomValidTest(rng, 50000))
	tests = append(tests, randomValidTest(rng, 200000))
	tests = append(tests, randomInvalidTest(rng, 20))
	tests = append(tests, randomInvalidTest(rng, 200))
	return tests
}

func newTestCase(n int, seq []int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range seq {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	b.WriteByte('\n')
	return testCase{input: b.String(), n: n}
}

func randomValidTest(rng *rand.Rand, n int) testCase {
	if n < 1 {
		n = 1
	}
	if n == 1 {
		return newTestCase(1, []int{0})
	}
	seqLen := 2*n - 1
	seq := make([]int, seqLen)
	center := rng.Intn(n) + 1
	used := make([]bool, n+1)
	used[center] = true
	cur := 1
	for v := 1; v <= n; v++ {
		if used[v] {
			continue
		}
		if cur >= seqLen {
			break
		}
		seq[cur] = v
		used[v] = true
		cur += 2
	}
	for i := 0; i < seqLen; i++ {
		if i%2 == 0 {
			seq[i] = center
		} else if seq[i] == 0 {
			seq[i] = randomUnused(rng, used, center)
		}
	}
	mask := make([]int, seqLen)
	for i := range seq {
		if rng.Intn(3) == 0 {
			mask[i] = seq[i]
		} else {
			mask[i] = 0
		}
	}
	return newTestCase(n, mask)
}

func randomUnused(rng *rand.Rand, used []bool, center int) int {
	n := len(used) - 1
	for {
		val := rng.Intn(n) + 1
		if val == center {
			continue
		}
		if !used[val] {
			used[val] = true
			return val
		}
	}
}

func randomInvalidTest(rng *rand.Rand, n int) testCase {
	if n < 2 {
		n = 2
	}
	seq := make([]int, 2*n-1)
	for i := range seq {
		seq[i] = rng.Intn(n + 1) // allow zeros
	}
	return newTestCase(n, seq)
}
