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

const refSource = "2048D.go"

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
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
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
			if len(candVals[i]) != len(refVals[i]) {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: expected %d numbers got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, len(refVals[i]), len(candVals[i]), tc.input, refOut, candOut)
				os.Exit(1)
			}
			for j := range refVals[i] {
				if candVals[i][j] != refVals[i][j] {
					fmt.Fprintf(os.Stderr, "test %d (%s) case %d position %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
						idx+1, tc.name, i+1, j+1, refVals[i][j], candVals[i][j], tc.input, refOut, candOut)
					os.Exit(1)
				}
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
	tmpDir, err := os.MkdirTemp("", "oracle-2048D-")
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

func parseOutput(output string, cases int) ([][]int64, error) {
	lines := strings.Split(output, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		filtered = append(filtered, line)
	}
	if len(filtered) != cases {
		return nil, fmt.Errorf("expected %d lines, got %d", cases, len(filtered))
	}
	res := make([][]int64, cases)
	for i := 0; i < cases; i++ {
		fields := strings.Fields(filtered[i])
		if len(fields) == 0 {
			return nil, fmt.Errorf("empty line %d", i+1)
		}
		vals := make([]int64, len(fields))
		for j, tok := range fields {
			v, err := strconv.ParseInt(tok, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q on line %d", tok, i+1)
			}
			vals[j] = v
		}
		res[i] = vals
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("sample", []manualCase{
			{
				n: 4, m: 4,
				a: []int{4, 3, 7, 5},
				b: []int{2, 5, 4, 6},
			},
			{
				n: 5, m: 5,
				a: []int{5, 0, 4, 8, 6},
				b: []int{1, 3, 9, 2, 7},
			},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

type manualCase struct {
	n, m int
	a    []int
	b    []int
}

func newManualTest(name string, cases []manualCase) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.m))
		for i, val := range c.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		for i, val := range c.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    name,
		input:   sb.String(),
		outputs: len(cases),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Intn(10)
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Intn(10)
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx+1),
		input:   sb.String(),
		outputs: t,
	}
}
