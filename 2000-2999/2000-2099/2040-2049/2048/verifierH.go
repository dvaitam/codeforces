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

const referenceSolutionRel = "2000-2999/2000-2099/2040-2049/2048/2048H.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2048H.go")
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
	s    []string
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(len(tc.s)*32 + 16)
	sb.WriteString(strconv.Itoa(len(tc.s)))
	sb.WriteByte('\n')
	for _, str := range tc.s {
		sb.WriteString(str)
		sb.WriteByte('\n')
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
	tmpDir, err := os.MkdirTemp("", "2048H-ref-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2048H")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, tok := range fields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		v %= 998244353
		if v < 0 {
			v += 998244353
		}
		res[i] = v
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample_single", s: []string{"11001000110111001100"}},
		{name: "single_one", s: []string{"1"}},
		{name: "single_zero", s: []string{"0"}},
		{name: "all_zero_short", s: []string{"0000"}},
		{name: "all_one_short", s: []string{"11111"}},
		{name: "alternating_even", s: []string{"01010101"}},
		{name: "alternating_odd", s: []string{"1010101"}},
		{name: "mixed_cases", s: []string{"0011", "1100", "010011"}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 160; i++ {
		t := rng.Intn(5) + 1
		arr := make([]string, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(80) + 1
			var sb strings.Builder
			sb.Grow(n)
			for k := 0; k < n; k++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			arr[j] = sb.String()
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			s:    arr,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, len(tc.s))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
		vals, err := parseOutputs(out, len(tc.s))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}

		for i := range vals {
			if vals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) mismatch at case %d: expected %d got %d\ninput:\n%sreference:\n%s\ncandidate:\n%s",
					idx+1, tc.name, i+1, refVals[i], vals[i], input, refOut, out)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed.")
}
