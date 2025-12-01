package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2030B.go"

type testCase struct {
	name string
	ns   []int
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
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		if err := validateOutput(tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed validation on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if err := validateOutput(tc, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2030B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2030B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.ns)))
	for _, n := range tc.ns {
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return sb.String()
}

func validateOutput(tc testCase, output string) error {
	lines := splitNonEmptyLines(output)
	if len(lines) != len(tc.ns) {
		return fmt.Errorf("expected %d output lines, got %d", len(tc.ns), len(lines))
	}
	for i, n := range tc.ns {
		if err := checkString(n, lines[i]); err != nil {
			return fmt.Errorf("case %d (n=%d): %v", i+1, n, err)
		}
	}
	return nil
}

func splitNonEmptyLines(out string) []string {
	raw := strings.Split(out, "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

func checkString(n int, s string) error {
	if len(s) != n {
		return fmt.Errorf("expected length %d, got %d", n, len(s))
	}
	zeros := 0
	for idx, ch := range s {
		if ch != '0' && ch != '1' {
			return fmt.Errorf("invalid character %q at position %d", ch, idx+1)
		}
		if ch == '0' {
			zeros++
		}
	}
	if n == 1 {
		return nil
	}
	if zeros != n-1 {
		ones := n - zeros
		return fmt.Errorf("must contain exactly one '1' (got zeros=%d ones=%d)", zeros, ones)
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "sample", ns: []int{1, 2, 3}},
		{name: "single_small", ns: []int{1}},
		{name: "single_large", ns: []int{200000}},
		{name: "mixed_small", ns: []int{2, 5, 10, 1, 4}},
	}

	rng := rand.New(rand.NewSource(987654321))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	target := rng.Intn(25) + 1
	remaining := 200000
	ns := make([]int, 0, target)
	for len(ns) < target && remaining > 0 {
		var maxVal int
		if rng.Intn(5) == 0 {
			maxVal = remaining
		} else {
			if remaining > 1000 {
				maxVal = 1000
			} else {
				maxVal = remaining
			}
		}
		if maxVal == 0 {
			break
		}
		val := rng.Intn(maxVal) + 1
		ns = append(ns, val)
		remaining -= val
	}
	if len(ns) == 0 {
		ns = append(ns, 1)
	}
	return testCase{
		name: fmt.Sprintf("random_%d", idx+1),
		ns:   ns,
	}
}
