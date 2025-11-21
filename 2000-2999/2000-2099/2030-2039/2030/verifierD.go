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

const refSource = "2030D.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate_binary")
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
		refVals, err := parseOutput(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if !strings.EqualFold(candVals[i], refVals[i]) {
				fmt.Fprintf(os.Stderr, "test %d (%s) mismatch on response %d: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
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
	tmpDir, err := os.MkdirTemp("", "oracle-2030D-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleD")
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

func parseOutput(output string, expected int) ([]string, error) {
	lines := strings.Fields(output)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(lines))
	}
	res := make([]string, expected)
	for i, ans := range lines {
		switch strings.ToUpper(ans) {
		case "YES", "NO":
			res[i] = strings.ToUpper(ans)
		default:
			return nil, fmt.Errorf("invalid answer %q, expected YES or NO", ans)
		}
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("sample_from_statement", 5, 3,
			[]int{1, 4, 2, 5, 3},
			"RLRLL",
			[]int{2, 4, 5},
		),
		newManualTest("already_sorted", 4, 2,
			[]int{1, 2, 3, 4},
			"RRLL",
			[]int{2, 3},
		),
		newManualTest("all_left_then_right", 6, 3,
			[]int{6, 5, 4, 3, 2, 1},
			"RRRRLL",
			[]int{3, 4, 5},
		),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, n, q int, p []int, s string, queries []int) testCase {
	if len(p) != n {
		panic("perm length mismatch")
	}
	if len(s) != n {
		panic("string length mismatch")
	}
	if len(queries) != q {
		panic("query count mismatch")
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(s)
	sb.WriteByte('\n')
	for _, qi := range queries {
		sb.WriteString(strconv.Itoa(qi))
		sb.WriteByte('\n')
	}
	return testCase{
		name:    name,
		input:   sb.String(),
		outputs: q,
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')

	n := rng.Intn(20) + 3
	q := rng.Intn(20) + 1
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	perm := randPerm(rng, n)
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	str := make([]byte, n)
	str[0] = 'R'
	str[n-1] = 'L'
	for i := 1; i < n-1; i++ {
		if rng.Intn(2) == 0 {
			str[i] = 'L'
		} else {
			str[i] = 'R'
		}
	}
	sb.Write(str)
	sb.WriteByte('\n')

	for i := 0; i < q; i++ {
		pos := rng.Intn(n-2) + 2
		sb.WriteString(strconv.Itoa(pos))
		sb.WriteByte('\n')
	}

	return testCase{
		name:    fmt.Sprintf("random_%d", idx+1),
		input:   sb.String(),
		outputs: q,
	}
}

func randPerm(rng *rand.Rand, n int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}
