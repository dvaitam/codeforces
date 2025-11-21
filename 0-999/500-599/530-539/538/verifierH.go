package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "0-999/500-599/530-539/538/538H.go"

type instance struct {
	t, T  int64
	n, m  int
	l, r  []int64
	edges [][2]int
}

type testCase struct {
	input string
	inst  instance
}

type candidateResult struct {
	possible   bool
	n1, n2     int64
	assignment string
}

func main() {
	if len(os.Args) != 2 {
		fatal("usage: go run verifierH.go /path/to/candidate")
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fatal("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fatal("reference failed on test %d: %v", idx+1, err)
		}
		refPossible, err := parsePossible(refOut)
		if err != nil {
			fatal("could not parse reference output on test %d: %v", idx+1, err)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fatal("candidate crashed on test %d: %v", idx+1, err)
		}
		candRes, err := parseCandidateOutput(candOut, tc.inst.n)
		if err != nil {
			fatal("failed to parse candidate output on test %d: %v", idx+1, err)
		}

		if refPossible {
			if !candRes.possible {
				fatal("candidate claims IMPOSSIBLE on solvable test %d", idx+1)
			}
			if err := validateSolution(tc.inst, candRes); err != nil {
				fatal("invalid solution on test %d: %v", idx+1, err)
			}
		} else {
			if candRes.possible {
				fatal("candidate outputs a solution on impossible test %d", idx+1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "538H-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runCmd(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runCmd(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCmd(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return stdout.String(), fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
		}
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func parsePossible(out string) (bool, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return false, fmt.Errorf("empty output")
	}
	verdict := strings.ToUpper(fields[0])
	switch verdict {
	case "POSSIBLE":
		return true, nil
	case "IMPOSSIBLE":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected verdict %q", fields[0])
	}
}

func parseCandidateOutput(out string, n int) (candidateResult, error) {
	reader := strings.NewReader(out)
	var verdict string
	if _, err := fmt.Fscan(reader, &verdict); err != nil {
		return candidateResult{}, fmt.Errorf("failed to read verdict: %w", err)
	}
	verdictUpper := strings.ToUpper(verdict)
	if verdictUpper == "IMPOSSIBLE" {
		return candidateResult{possible: false}, nil
	}
	if verdictUpper != "POSSIBLE" {
		return candidateResult{}, fmt.Errorf("unknown verdict %q", verdict)
	}
	var n1, n2 int64
	if _, err := fmt.Fscan(reader, &n1, &n2); err != nil {
		return candidateResult{}, fmt.Errorf("failed to read group sizes: %w", err)
	}
	var assignment string
	if _, err := fmt.Fscan(reader, &assignment); err != nil {
		return candidateResult{}, fmt.Errorf("failed to read assignment string: %w", err)
	}
	if len(assignment) != n {
		return candidateResult{}, fmt.Errorf("assignment length %d does not match n=%d", len(assignment), n)
	}
	for i, ch := range assignment {
		if ch != '1' && ch != '2' {
			return candidateResult{}, fmt.Errorf("invalid character %q at position %d", ch, i)
		}
	}
	return candidateResult{possible: true, n1: n1, n2: n2, assignment: assignment}, nil
}

func validateSolution(inst instance, res candidateResult) error {
	if res.n1 < 0 || res.n2 < 0 {
		return fmt.Errorf("group sizes must be non-negative")
	}
	total := res.n1 + res.n2
	if total < inst.t || total > inst.T {
		return fmt.Errorf("total students %d outside [%d,%d]", total, inst.t, inst.T)
	}
	if len(res.assignment) != inst.n {
		return fmt.Errorf("assignment length mismatch")
	}
	for i := 0; i < inst.n; i++ {
		group := res.assignment[i]
		size := res.n1
		if group == '2' {
			size = res.n2
		}
		if size < inst.l[i] || size > inst.r[i] {
			return fmt.Errorf("teacher %d violates interval [%d,%d] with group size %d", i+1, inst.l[i], inst.r[i], size)
		}
	}
	for _, e := range inst.edges {
		if res.assignment[e[0]] == res.assignment[e[1]] {
			return fmt.Errorf("teachers %d and %d must be in different groups", e[0]+1, e[1]+1)
		}
	}
	return nil
}

func parseInstance(input string) (instance, error) {
	reader := strings.NewReader(input)
	var inst instance
	if _, err := fmt.Fscan(reader, &inst.t, &inst.T); err != nil {
		return inst, err
	}
	if _, err := fmt.Fscan(reader, &inst.n, &inst.m); err != nil {
		return inst, err
	}
	inst.l = make([]int64, inst.n)
	inst.r = make([]int64, inst.n)
	for i := 0; i < inst.n; i++ {
		if _, err := fmt.Fscan(reader, &inst.l[i], &inst.r[i]); err != nil {
			return inst, err
		}
	}
	inst.edges = make([][2]int, inst.m)
	for i := 0; i < inst.m; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return inst, err
		}
		inst.edges[i] = [2]int{u - 1, v - 1}
	}
	return inst, nil
}

func buildTests() []testCase {
	inputs := []string{
		"10 20\n3 0\n3 6\n4 9\n16 25\n",
		"1 10\n3 3\n0 10\n0 10\n0 10\n1 2\n1 3\n2 3\n",
		"1 5\n2 1\n0 5\n0 5\n1 2\n",
		"5 5\n2 0\n5 5\n0 2\n",
		"5 5\n2 0\n3 3\n3 3\n",
		"4 9\n4 3\n1 5\n2 6\n1 4\n3 7\n1 2\n2 3\n3 4\n",
	}
	tests := make([]testCase, 0, len(inputs))
	for _, in := range inputs {
		inst, err := parseInstance(in)
		if err != nil {
			panic(fmt.Sprintf("failed to parse built-in test: %v", err))
		}
		tests = append(tests, testCase{input: in, inst: inst})
	}
	return tests
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
