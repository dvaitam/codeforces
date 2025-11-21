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

type testCase struct {
	input      string
	caseCount  int
	expectedNs []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refTokens, err := parseOutputs(refOut, tc.caseCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotTokens, err := parseOutputs(gotOut, tc.caseCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.caseCount; caseIdx++ {
			if refTokens[caseIdx] != gotTokens[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %s got %s\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, caseIdx+1, refTokens[caseIdx], gotTokens[caseIdx], tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2172E_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2172E.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, caseCount int) ([]string, error) {
	lines := strings.Fields(strings.TrimSpace(out))
	if len(lines) != caseCount {
		return nil, fmt.Errorf("expected %d tokens got %d", caseCount, len(lines))
	}
	for _, tok := range lines {
		if err := validateToken(tok); err != nil {
			return nil, err
		}
	}
	return lines, nil
}

func validateToken(tok string) error {
	if !strings.Contains(tok, "A") || !strings.Contains(tok, "B") {
		return fmt.Errorf("invalid token %q", tok)
	}
	parts := strings.Split(tok, "A")
	if len(parts) != 2 {
		return fmt.Errorf("invalid token %q", tok)
	}
	left := parts[0]
	if left == "" {
		return fmt.Errorf("invalid A value in %q", tok)
	}
	if _, err := strconv.Atoi(left); err != nil {
		return fmt.Errorf("invalid A integer in %q", tok)
	}
	rightParts := strings.Split(parts[1], "B")
	if len(rightParts) != 2 || rightParts[1] != "" {
		return fmt.Errorf("invalid B format in %q", tok)
	}
	if _, err := strconv.Atoi(rightParts[0]); err != nil {
		return fmt.Errorf("invalid B integer in %q", tok)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 200)...)
	return tests
}

func manualTests() []testCase {
	return []testCase{
		makeTestCase([][3]int{
			{12, 1, 1},
			{12, 1, 2},
			{123, 1, 2},
			{123, 2, 5},
			{1234, 1, 24},
			{1234, 15, 9},
			{1234, 1, 1},
		}),
	}
}

func randomTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	options := []int{12, 123, 1234}
	length := map[int]int{12: 2, 123: 3, 1234: 4}
	fact := map[int]int{2: 2, 3: 6, 4: 24}
	for b := 0; b < batches; b++ {
		caseCount := rng.Intn(10) + 1
		entries := make([][3]int, caseCount)
		for i := 0; i < caseCount; i++ {
			n := options[rng.Intn(len(options))]
			total := fact[length[n]]
			j := rng.Intn(total) + 1
			k := rng.Intn(total) + 1
			entries[i] = [3]int{n, j, k}
		}
		tests = append(tests, makeTestCase(entries))
	}
	return tests
}

func makeTestCase(entries [][3]int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(entries)))
	sb.WriteByte('\n')
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return testCase{
		input:     sb.String(),
		caseCount: len(entries),
	}
}
