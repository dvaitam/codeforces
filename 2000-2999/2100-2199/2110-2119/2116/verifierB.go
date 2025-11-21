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

const refSource = "2000-2999/2100-2199/2110-2119/2116/2116B.go"

type testCase struct {
	name    string
	input   string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	dir, err := os.MkdirTemp("", "cf-2116B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2116B.bin")
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

func parseAnswers(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func equalSlices(a, b []int64) bool {
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
		newTestCase("sample", "3\n3\n0 2 1\n1 2 0\n5\n0 1 2 3 4\n4 3 2 1 0\n10\n5 8 9 3 4 0 2 7 1 6\n9 5 1 4 0 3 2 8 7 6\n"),
		buildSingle("n1", []int{0}, []int{0}),
		buildSingle("n2_small", []int{0, 1}, []int{1, 0}),
		buildSingle("n4_perm", []int{2, 0, 3, 1}, []int{1, 3, 0, 2}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func buildSingle(name string, p, q []int) testCase {
	if len(p) != len(q) {
		panic("p and q length mismatch")
	}
	n := len(p)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range q {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String(), answers: n}
}

func newTestCase(name, input string) testCase {
	cnt, err := countAnswers(input)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test %s: %v", name, err))
	}
	return testCase{name: name, input: input, answers: cnt}
}

func countAnswers(input string) (int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, fmt.Errorf("failed to read t: %v", err)
	}
	if t <= 0 {
		return 0, fmt.Errorf("non-positive t: %d", t)
	}
	total := 0
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return 0, fmt.Errorf("failed to read n: %v", err)
		}
		for i := 0; i < n; i++ {
			var tmp int
			if _, err := fmt.Fscan(reader, &tmp); err != nil {
				return 0, fmt.Errorf("failed to read p[%d]: %v", i, err)
			}
		}
		for i := 0; i < n; i++ {
			var tmp int
			if _, err := fmt.Fscan(reader, &tmp); err != nil {
				return 0, fmt.Errorf("failed to read q[%d]: %v", i, err)
			}
		}
		total += n
	}
	return total, nil
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	totalAns := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(12) + 1
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		p := rng.Perm(n)
		q := rng.Perm(n)
		for j, v := range p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for j, v := range q {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		totalAns += n
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		answers: totalAns,
	}
}
