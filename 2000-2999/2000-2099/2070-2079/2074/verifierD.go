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

const refSource = "2074D.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate_binary_or_go_file")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin := candidate
	candCleanup := func() {}
	if strings.HasSuffix(candidate, ".go") {
		candBin, candCleanup, err = buildBinary(candidate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
			os.Exit(1)
		}
	}
	defer candCleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at answer %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildBinary(src string) (string, func(), error) {
	abs, err := filepath.Abs(src)
	if err != nil {
		return "", nil, err
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier path")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "2074D-bin-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "bin")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
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

func parseOutputs(output string, expected int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("single_circle", []circleCase{
			{n: 1, x: []int64{0}, r: []int{2}},
		}),
		newManualTest("two_overlap", []circleCase{
			{n: 2, x: []int64{0, 0}, r: []int{1, 2}},
		}),
		newManualTest("sample", parseSampleInput()),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

type circleCase struct {
	n int
	x []int64
	r []int
}

func newManualTest(name string, cases []circleCase) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, cs := range cases {
		m := 0
		for _, v := range cs.r {
			m += v
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, m))
		for i, v := range cs.x {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range cs.r {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String(), outputs: len(cases)}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	cases := make([]circleCase, t)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(60) + n // ensure m >= n
		radii := randomRadii(rng, n, m)
		xs := make([]int64, n)
		for j := 0; j < n; j++ {
			xs[j] = int64(rng.Intn(201) - 100)
		}
		cases[i] = circleCase{n: n, x: xs, r: radii}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j, v := range xs {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for j, v := range radii {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx+1),
		input:   sb.String(),
		outputs: t,
	}
}

func randomRadii(rng *rand.Rand, n int, m int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = 1
	}
	remaining := m - n
	for remaining > 0 {
		idx := rng.Intn(n)
		add := rng.Intn(remaining + 1)
		res[idx] += add
		remaining -= add
	}
	rng.Shuffle(n, func(i, j int) { res[i], res[j] = res[j], res[i] })
	return res
}

func parseSampleInput() []circleCase {
	// Sample from the statement.
	return []circleCase{
		{n: 2, x: []int64{0, 0}, r: []int{1, 2}},
		{n: 2, x: []int64{0, 2}, r: []int{1, 2}},
		{n: 3, x: []int64{0, 2, 5}, r: []int{1, 1, 1}},
		{n: 4, x: []int64{0, 5, 10, 15}, r: []int{2, 2, 2, 2}},
	}
}
