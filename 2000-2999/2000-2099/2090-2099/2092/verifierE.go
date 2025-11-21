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

const refSource = "2000-2999/2000-2099/2090-2099/2092/2092E.go"

type testCase struct {
	name    string
	input   string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	dir, err := os.MkdirTemp("", "cf-2092E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2092E.bin")
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
		newTestCase("sample", "2\n3 3 6\n1 1 0\n1 2 1\n1 3 0\n3 1 1\n3 2 0\n3 3 1\n3 4 12\n1 1 0\n1 2 1\n1 3 0\n1 4 1\n2 1 1\n2 2 0\n2 3 1\n2 4 0\n3 1 0\n3 2 1\n3 3 0\n3 4 1\n"),
		buildSingleCase("all_green_small", 3, 3, 1),
		buildSingleCase("single_colored", 4, 5, 1),
		buildSingleCase("large_board", 1000000000, 1000000000, 2),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func buildSingleCase(name string, n, m int64, k int) testCase {
	if k < 0 {
		panic("negative k")
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < k; i++ {
		x := int64(1 + i%3)
		y := int64(1 + (i/3)%3)
		c := i & 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, c))
	}
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
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	totalK := 0
	for i := 0; i < t; i++ {
		n := int64(rng.Intn(18) + 3)
		m := int64(rng.Intn(18) + 3)
		maxK := int(n * m) // small enough here
		k := rng.Intn(minInt(maxK, 30)) + 1
		totalK += k
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		used := make(map[int]struct{}, k)
		for j := 0; j < k; j++ {
			for {
				x := rng.Intn(int(n)) + 1
				y := rng.Intn(int(m)) + 1
				key := x*1000 + y
				if _, ok := used[key]; ok {
					continue
				}
				used[key] = struct{}{}
				c := rng.Intn(2)
				sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, c))
				break
			}
		}
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		answers: t,
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
