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

const refSource = "2040C.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate_binary")
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

		for caseIdx := 0; caseIdx < tc.outputs; caseIdx++ {
			ref := refVals[caseIdx]
			got := candVals[caseIdx]
			if ref.validCount != got.validCount {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: mismatch answer count expectation\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, caseIdx+1, tc.input, refOut, candOut)
				os.Exit(1)
			}
			if ref.validCount == 0 {
				continue
			}
			if len(ref.perm) != len(got.perm) {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: length mismatch\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, caseIdx+1, tc.input, refOut, candOut)
				os.Exit(1)
			}
			for i := range ref.perm {
				if ref.perm[i] != got.perm[i] {
					fmt.Fprintf(os.Stderr, "test %d (%s) case %d: permutation mismatch\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
						idx+1, tc.name, caseIdx+1, tc.input, refOut, candOut)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

type parsedAnswer struct {
	validCount int // 0 if -1, 1 if permutation printed
	perm       []int
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2040C-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleC")
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

func parseOutput(output string, cases int) ([]parsedAnswer, error) {
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
	res := make([]parsedAnswer, cases)
	for i := 0; i < cases; i++ {
		line := filtered[i]
		if line == "" {
			return nil, fmt.Errorf("empty line %d", i+1)
		}
		if line == "-1" {
			res[i] = parsedAnswer{}
			continue
		}
		fields := strings.Fields(line)
		perm := make([]int, len(fields))
		used := make(map[int]bool, len(fields))
		for idx, tok := range fields {
			val, err := strconv.Atoi(tok)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q on line %d", tok, i+1)
			}
			if val < 1 || val > len(fields) {
				return nil, fmt.Errorf("value %d out of range on line %d", val, i+1)
			}
			if used[val] {
				return nil, fmt.Errorf("duplicate value %d on line %d", val, i+1)
			}
			used[val] = true
			perm[idx] = val
		}
		res[i] = parsedAnswer{validCount: 1, perm: perm}
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("small_cases", []testQuery{
			{n: 3, k: 1},
			{n: 3, k: 2},
			{n: 3, k: 3},
			{n: 3, k: 4},
		}),
		newManualTest("edge_k", []testQuery{
			{n: 1, k: 1},
			{n: 1, k: 2},
			{n: 2, k: 1},
			{n: 2, k: 4},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

type testQuery struct {
	n int
	k int64
}

func newManualTest(name string, queries []testQuery) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.n, q.k))
	}
	return testCase{
		name:    name,
		input:   sb.String(),
		outputs: len(queries),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(15) + 1
		if rng.Intn(5) == 0 {
			n = rng.Intn(50) + 1
		}
		var k int64
		if rng.Intn(3) == 0 {
			k = int64(rng.Intn(20) + 1)
		} else {
			k = int64(1 + rng.Intn(1_000_000_000))
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx+1),
		input:   sb.String(),
		outputs: t,
	}
}
