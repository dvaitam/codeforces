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

const referenceSolutionRel = "2000-2999/2000-2099/2030-2039/2038/2038G.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2038G.go")
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
	strs []string
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(len(tc.strs)*64 + 16)
	sb.WriteString(strconv.Itoa(len(tc.strs)))
	sb.WriteByte('\n')
	for _, s := range tc.strs {
		sb.WriteString(fmt.Sprintf("%d\n%s\n", len(s), s))
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
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2038G-ref-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2038G")
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

func parseGuesses(out string, tc testCase) error {
	fields := strings.Fields(out)
	expected := len(tc.strs) * 2
	if len(fields) != expected {
		return fmt.Errorf("expected %d tokens (index+value per test), got %d", expected, len(fields))
	}
	pos := 0
	for caseIdx, s := range tc.strs {
		n := len(s)
		idxVal := fields[pos]
		pos++
		valTok := fields[pos]
		pos++

		idx, err := strconv.Atoi(idxVal)
		if err != nil {
			return fmt.Errorf("test %d: invalid index %q", caseIdx+1, idxVal)
		}
		if idx < 1 || idx > n {
			return fmt.Errorf("test %d: index %d out of range 1..%d", caseIdx+1, idx, n)
		}
		if valTok != "0" && valTok != "1" {
			return fmt.Errorf("test %d: value must be 0 or 1, got %q", caseIdx+1, valTok)
		}
		if byte(valTok[0]) != s[idx-1] {
			return fmt.Errorf("test %d: guessed %s at position %d but string has %c", caseIdx+1, valTok, idx, s[idx-1])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "single_case_all_zero", strs: []string{"0000"}},
		{name: "single_case_all_one", strs: []string{"11111"}},
		{name: "two_cases_mixed", strs: []string{"01", "10"}},
		{name: "alternating_even", strs: []string{"010101"}},
		{name: "alternating_odd", strs: []string{"1010101"}},
		{name: "small_variants", strs: []string{"01", "11", "00", "101", "0001"}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 150; i++ {
		t := rng.Intn(6) + 1
		strs := make([]string, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(49) + 2
			var sb strings.Builder
			sb.Grow(n)
			for k := 0; k < n; k++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			strs[j] = sb.String()
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			strs: strs,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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

		// Ensure reference solution agrees with the input format.
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		if err := parseGuesses(refOut, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
		if err := parseGuesses(out, tc); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
