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

const refSource = "2000-2999/2000-2099/2090-2099/2097/2097A.go"

type testCase struct {
	name    string
	input   string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if !equalSlices(refAns, candAns) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\ninput:\n%sreference:\n%s\ncandidate:\n%s", idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2097A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2097A.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(output string, expected int) ([]string, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]string, expected)
	for i, token := range fields {
		res[i] = strings.ToLower(token)
	}
	return res, nil
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample", "5\n4\n1 1 1 1\n3\n2 2 2\n5\n2 4 3 2 4\n8\n6 3 1 1 5 1 2 6\n1\n1000000000\n"),
		buildSingle("n1", 1),
		buildSingle("n3", 3),
		buildSingle("n4", 4),
		buildSingle("nLarge", 100000),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func buildSingle(name string, n int) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(i + 1))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String(), answers: 1}
}

func newTestCase(name, input string) testCase {
	cnt, err := countCases(input)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test %s: %v", name, err))
	}
	return testCase{name: name, input: input, answers: cnt}
}

func countCases(input string) (int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, fmt.Errorf("failed to read t: %v", err)
	}
	if t <= 0 {
		return 0, fmt.Errorf("non-positive t: %d", t)
	}
	return t, nil
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Int63n(1_000_000_000) + 1
			sb.WriteString(strconv.FormatInt(val, 10))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		answers: t,
	}
}
