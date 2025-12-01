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

const refSource2001B = "./2001B.go"

type testCase struct {
	name  string
	cases []int
	input string
}

type caseResult struct {
	impossible bool
	perm       []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		if err := validateOutput(tc.cases, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := validateOutput(tc.cases, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2001B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2001B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2001B)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validateOutput(ns []int, output string) error {
	results, err := parseOutput(ns, output)
	if err != nil {
		return err
	}
	for idx, n := range ns {
		res := results[idx]
		if hasSolution(n) {
			if res.impossible {
				return fmt.Errorf("case %d (n=%d) should have a solution but output -1", idx+1, n)
			}
			if err := validatePermutation(res.perm, n); err != nil {
				return fmt.Errorf("case %d (n=%d): %v", idx+1, n, err)
			}
			if err := validateCarriageParity(res.perm); err != nil {
				return fmt.Errorf("case %d (n=%d): %v", idx+1, n, err)
			}
		} else {
			if !res.impossible {
				return fmt.Errorf("case %d (n=%d) should be impossible but permutation %v was printed", idx+1, n, res.perm)
			}
		}
	}
	return nil
}

func hasSolution(n int) bool {
	return n%2 == 1
}

func parseOutput(ns []int, out string) ([]caseResult, error) {
	tokens := strings.Fields(out)
	results := make([]caseResult, len(ns))
	idx := 0
	for caseIdx, n := range ns {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("case %d: missing output", caseIdx+1)
		}
		token := tokens[idx]
		idx++
		if token == "-1" {
			results[caseIdx] = caseResult{impossible: true}
			continue
		}
		perm := make([]int, n)
		val, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("case %d: invalid integer %q", caseIdx+1, token)
		}
		perm[0] = val
		for i := 1; i < n; i++ {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("case %d: expected %d numbers, got %d", caseIdx+1, n, i)
			}
			token = tokens[idx]
			idx++
			val, err = strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("case %d: invalid integer %q", caseIdx+1, token)
			}
			perm[i] = val
		}
		results[caseIdx] = caseResult{perm: perm}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra token %q in output", tokens[idx])
	}
	return results, nil
}

func validatePermutation(perm []int, n int) error {
	if len(perm) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(perm))
	}
	seen := make([]bool, n+1)
	for i, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d at position %d is outside [1,%d]", v, i+1, n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
	}
	return nil
}

func validateCarriageParity(perm []int) error {
	n := len(perm)
	if n == 0 {
		return fmt.Errorf("empty permutation")
	}
	pos := make([]int, n+1)
	for idx, v := range perm {
		pos[v] = idx
	}
	asc, desc := 0, 0
	for v := 1; v < n; v++ {
		if pos[v+1] > pos[v] {
			asc++
		} else {
			desc++
		}
	}
	if asc != desc {
		return fmt.Errorf("needs %d resets from right and %d from left", asc, desc)
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("n1", []int{1}),
		newTestCase("n2", []int{2}),
		newTestCase("n3", []int{3}),
		newTestCase("mixed_small", []int{1, 3, 5}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 120; i++ {
		caseCnt := rng.Intn(4) + 1
		ns := make([]int, caseCnt)
		for j := 0; j < caseCnt; j++ {
			switch rng.Intn(5) {
			case 0:
				ns[j] = 1
			case 1:
				ns[j] = 2
			default:
				n := rng.Intn(100) + 1
				if n%2 == 0 {
					n++
				}
				ns[j] = n
			}
		}
		tests = append(tests, newTestCase(fmt.Sprintf("random_odd_%d", i+1), ns))
	}

	return tests
}

func newTestCase(name string, ns []int) testCase {
	return testCase{
		name:  name,
		cases: ns,
		input: formatInput(ns),
	}
}

func formatInput(ns []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ns)))
	for _, n := range ns {
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return sb.String()
}
