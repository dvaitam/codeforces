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
)

const (
	refSource = "./338D.go"
	maxNM     = int64(1_000_000_000_000)
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expectOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expectVerdict, err := parseVerdict(expectOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, expectOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVerdict, err := parseVerdict(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}

		if gotVerdict != expectVerdict {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %s, got %s\ninput:\n%s\n", idx+1, tc.name, expectVerdict, gotVerdict, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "338D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseVerdict(out string) (string, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return "", fmt.Errorf("expected single token, got %q", out)
	}
	v := strings.ToUpper(fields[0])
	if v != "YES" && v != "NO" {
		return "", fmt.Errorf("verdict must be YES or NO, got %q", fields[0])
	}
	return v, nil
}

func buildPositiveByRow(name string, n, m, row, start int64, k int) testCase {
	seq := make([]int64, k)
	for i := 0; i < k; i++ {
		seq[i] = gcd(row, start+int64(i))
	}
	return testCase{name: name, input: formatInput(n, m, seq)}
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(3380338))
	var tests []testCase

	tests = append(tests,
		buildPositiveByRow("small-positive", 10, 10, 10, 5, 5),
		testCase{name: "too-short-width", input: formatInput(5, 4, []int64{1, 2, 3, 4, 5})},
	)

	tests = append(tests, buildPositiveByRow("fixed-positive-1", 100, 200, 60, 50, 8))
	tests = append(tests, buildPositiveByRow("fixed-positive-2", 5000, 8000, 4200, 1000, 20))
	tests = append(tests, testCase{name: "value-exceeds-n", input: formatInput(50, 200, []int64{1, 2, 3, 51})})

	tests = append(tests, largePositiveCase())
	tests = append(tests, largeNegativeCase())

	for i := 0; i < 25; i++ {
		tests = append(tests, randomPositiveCase(rng, fmt.Sprintf("rand-positive-%d", i+1), 300))
	}
	for i := 0; i < 25; i++ {
		tests = append(tests, randomNegativeCase(rng, fmt.Sprintf("rand-negative-%d", i+1), 300))
	}

	return tests
}

func largePositiveCase() testCase {
	k := 10000
	n := int64(1_000_000_000_000)
	m := int64(k) + 5000
	row := int64(999_999_937)
	start := int64(2000)
	return buildPositiveByRow("large-positive-10k", n, m, row, start, k)
}

func largeNegativeCase() testCase {
	k := 10000
	n := int64(1_000_000)
	m := int64(9999)
	seq := make([]int64, k)
	for i := range seq {
		seq[i] = 1 + int64(i%5)
	}
	return testCase{name: "large-negative-too-wide", input: formatInput(n, m, seq)}
}

func randomPositiveCase(rng *rand.Rand, name string, maxK int) testCase {
	k := rng.Intn(maxK) + 1
	n := rng.Int63n(maxNM) + 1
	mMin := int64(k)
	m := mMin + rng.Int63n(maxNM-mMin+1)
	row := rng.Int63n(n) + 1
	startMax := m - int64(k) + 1
	start := rng.Int63n(startMax) + 1
	seq := make([]int64, k)
	for i := 0; i < k; i++ {
		seq[i] = gcd(row, start+int64(i))
	}
	return testCase{name: name, input: formatInput(n, m, seq)}
}

func randomNegativeCase(rng *rand.Rand, name string, maxK int) testCase {
	k := rng.Intn(maxK) + 1
	n := rng.Int63n(maxNM) + 1
	m := rng.Int63n(maxNM) + 1
	seq := make([]int64, k)
	for i := 0; i < k; i++ {
		seq[i] = rng.Int63n(maxNM) + 1
	}

	switch rng.Intn(3) {
	case 0:
		if int64(k) <= m {
			m = int64(k) - 1
			if m < 1 {
				m = 1
			}
		}
	case 1:
		idx := rng.Intn(k)
		seq[idx] = n + rng.Int63n(1000) + 1
	}

	return testCase{name: name, input: formatInput(n, m, seq)}
}

func formatInput(n, m int64, seq []int64) string {
	var sb strings.Builder
	sb.Grow(32 + len(seq)*14)
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, len(seq))
	for i, v := range seq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}
