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

const referenceSolutionRel = "2000-2999/2000-2099/2070-2079/2075/2075F.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2075F.go")
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
	a    []int
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.Grow(len(cases) * 64)
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.a)))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
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
	tmpDir, err := os.MkdirTemp("", "2075F-ref-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2075F")
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

func parseOutputs(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	ans := make([]int, t)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", f, i+1)
		}
		ans[i] = val
	}
	return ans, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample1", a: []int{4}},
		{name: "sample2", a: []int{5, 6, 5, 4, 3, 2, 1}},
		{name: "sample3", a: []int{1, 1, 3, 4, 2, 3, 4}},
		{name: "sample4", a: []int{2, 3, 1, 1, 2, 4}},
		{name: "increasing", a: []int{1, 2, 3, 4, 5, 6}},
		{name: "decreasing", a: []int{6, 5, 4, 3, 2, 1}},
		{name: "all_equal", a: []int{7, 7, 7, 7}},
		{name: "small_mix", a: []int{1, 4, 2, 4, 7}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 220; i++ {
		n := rng.Intn(200) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(1_000_000) + 1
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			a:    arr,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		input := buildInput([]testCase{tc})

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutputs(refOut, 1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
		ans, err := parseOutputs(out, 1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}

		if ans[0] != refAns[0] {
			fmt.Fprintf(os.Stderr, "test %d (%s) mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refAns[0], ans[0], input, refOut, out)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
