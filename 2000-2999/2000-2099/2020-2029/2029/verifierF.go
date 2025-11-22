package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const referenceSolutionRel = "2000-2999/2000-2099/2020-2029/2029/2029F.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2029F.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	name string
	n    int
	s    string
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.Grow(len(cases) * 32)
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n%s\n", tc.n, tc.s))
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2029F-ref-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2029F")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutputs(out string, t int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	ans := make([]string, t)
	for i, f := range fields {
		l := strings.ToLower(f)
		if l == "yes" || l == "no" {
			ans[i] = l
			continue
		}
		return nil, fmt.Errorf("invalid answer token %q", f)
	}
	return ans, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample1", n: 5, s: "RRRRR"},
		{name: "sample2", n: 5, s: "RRRRB"},
		{name: "sample3", n: 5, s: "RBBRB"},
		{name: "sample4", n: 6, s: "RBRBRB"},
		{name: "sample5", n: 6, s: "RRBBRB"},
		{name: "sample6", n: 5, s: "RBRBR"},
		{name: "sample7", n: 12, s: "RRBRRBRRBRRB"},
		{name: "all_red_small", n: 3, s: "RRR"},
		{name: "single_blue", n: 7, s: "RRRRBRR"},
		{name: "even_alternate", n: 8, s: "RBRBRBRB"},
		{name: "odd_alternate", n: 9, s: "BRBRBRBRB"},
		{name: "double_rr_only", n: 7, s: "RRBRBRB"},
		{name: "double_bb_only", n: 7, s: "BBRBRBR"},
		{name: "isolated_target_even_gaps", n: 9, s: "RRRBRBBRR"},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 150; i++ {
		n := rng.Intn(120) + 3
		var sb strings.Builder
		sb.Grow(n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('R')
			} else {
				sb.WriteByte('B')
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			n:    n,
			s:    sb.String(),
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := buildInput([]testCase{tc})

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		refAns, err := parseOutputs(refOut, 1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid reference output on test %d (%s): %v\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		ans, err := parseOutputs(out, 1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\n%s", idx+1, tc.name, err, out)
			os.Exit(1)
		}

		if ans[0] != refAns[0] {
			fmt.Fprintf(os.Stderr, "test %d (%s): expected %s, got %s\ninput:\n%s", idx+1, tc.name, refAns[0], ans[0], input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
