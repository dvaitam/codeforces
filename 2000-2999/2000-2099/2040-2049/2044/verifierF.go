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

const refSource = "2044F.go"

type testCase struct {
	name         string
	input        string
	totalQueries int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate_binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, tc.totalQueries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, tc.totalQueries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for qIdx := 0; qIdx < tc.totalQueries; qIdx++ {
			if !strings.EqualFold(refAns[qIdx], candAns[qIdx]) {
				fmt.Fprintf(os.Stderr, "test %d (%s) mismatch at query %d: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, qIdx+1, refAns[qIdx], candAns[qIdx], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2044F-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(output string, count int) ([]string, error) {
	lines := strings.Split(output, "\n")
	ans := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch strings.ToUpper(line) {
		case "YES", "NO":
			ans = append(ans, strings.ToUpper(line))
		default:
			return nil, fmt.Errorf("invalid answer %q", line)
		}
	}
	if len(ans) != count {
		return nil, fmt.Errorf("expected %d answers, got %d", count, len(ans))
	}
	return ans, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("simple", []int64{1, 2}, []int64{3, 4}, []int{1, 2, 3}),
		newManualTest("zero_vectors", []int64{0, 0, 0}, []int64{0, 1}, []int{1, -1}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, a []int64, b []int64, queries []int) testCase {
	n := len(a)
	m := len(b)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, len(queries)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d\n", q))
	}
	return testCase{
		name:         name,
		input:        sb.String(),
		totalQueries: len(queries),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(8) + 1
	m := rng.Intn(8) + 1
	q := rng.Intn(8) + 1

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := rng.Intn(2*n+1) - n
		sb.WriteString(strconv.Itoa(val))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := rng.Intn(2*m+1) - m
		sb.WriteString(strconv.Itoa(val))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		val := rng.Intn(400001) - 200000
		if val == 0 {
			val = 1
		}
		sb.WriteString(fmt.Sprintf("%d\n", val))
	}
	return testCase{
		name:         fmt.Sprintf("random_%d", idx+1),
		input:        sb.String(),
		totalQueries: q,
	}
}
