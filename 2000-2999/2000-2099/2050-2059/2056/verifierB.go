package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const referenceSolutionRel = "2000-2999/2000-2099/2050-2059/2056/2056B.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2056B.go")
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
	mat  []string
}

func buildMatrix(perm []int) []string {
	n := len(perm)
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			grid[i][j] = '0'
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a, b := perm[i]-1, perm[j]-1
			if perm[i] < perm[j] {
				grid[a][b] = '1'
				grid[b][a] = '1'
			}
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(grid[i])
	}
	return res
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n*(tc.n+2) + 16)
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, row := range tc.mat {
		sb.WriteString(row)
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
	tmpDir, err := os.MkdirTemp("", "2056B-ref-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2056B")
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

func parsePermutation(out string, n int) ([]int, error) {
	reader := strings.NewReader(out)
	perm := make([]int, 0, n)
	for len(perm) < n {
		var v int
		if _, err := fmt.Fscan(reader, &v); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d integers, got %d", n, len(perm))
			}
			return nil, fmt.Errorf("failed to read integer %d: %v", len(perm)+1, err)
		}
		perm = append(perm, v)
	}
	var extra int
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extraneous output after %d integers", n)
	}
	return perm, nil
}

func validatePermutation(perm []int, tc testCase) error {
	if len(perm) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(perm))
	}
	seen := make([]bool, tc.n+1)
	for i, v := range perm {
		if v < 1 || v > tc.n {
			return fmt.Errorf("position %d: value %d out of range 1..%d", i+1, v, tc.n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
	}
	// verify adjacency consistency
	for i := 0; i < tc.n; i++ {
		for j := i + 1; j < tc.n; j++ {
			a, b := perm[i]-1, perm[j]-1
			edge := tc.mat[a][b] == '1'
			expect := perm[i] < perm[j]
			if edge != expect {
				return fmt.Errorf("pair (%d,%d) mismatch: edge=%v expect=%v", perm[i], perm[j], edge, expect)
			}
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		func() testCase {
			perm := []int{1}
			return testCase{name: "single", n: len(perm), mat: buildMatrix(perm)}
		}(),
		func() testCase {
			perm := []int{1, 3, 2, 4}
			return testCase{name: "small_mix", n: len(perm), mat: buildMatrix(perm)}
		}(),
		func() testCase {
			perm := []int{5, 4, 3, 2, 1}
			return testCase{name: "strict_desc", n: len(perm), mat: buildMatrix(perm)}
		}(),
		func() testCase {
			perm := []int{1, 2, 3, 4, 5, 6, 7}
			return testCase{name: "strict_inc", n: len(perm), mat: buildMatrix(perm)}
		}(),
		func() testCase {
			perm := []int{2, 5, 1, 3, 4}
			return testCase{name: "medium_mix", n: len(perm), mat: buildMatrix(perm)}
		}(),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 180; i++ {
		n := rng.Intn(80) + 1
		order := rng.Perm(n)
		perm := make([]int, n)
		for j, v := range order {
			perm[j] = v + 1
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			n:    n,
			mat:  buildMatrix(perm),
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		refPerm, err := parsePermutation(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}
		if err := validatePermutation(refPerm, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed validation on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
		perm, err := parsePermutation(out, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
		if err := validatePermutation(perm, tc); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed validation on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
