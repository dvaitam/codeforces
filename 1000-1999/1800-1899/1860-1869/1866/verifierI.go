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

type testCase struct {
	name  string
	input string
}

var verifierDir string

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		verifierDir = filepath.Dir(file)
	} else {
		verifierDir = "."
	}
}

func buildReference() (string, error) {
	outPath := filepath.Join(verifierDir, "ref1866I.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1866I.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	if !filepath.IsAbs(target) {
		if abs, err := filepath.Abs(target); err == nil {
			target = abs
		}
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
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

func parseWinner(out string) (string, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty output")
	}
	w := fields[0]
	if w != "Chaneka" && w != "Bhinneka" {
		return "", fmt.Errorf("invalid winner %q", w)
	}
	return w, nil
}

func verifyCase(candidate, reference string, tc testCase) error {
	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	expected, err := parseWinner(refOut)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate error: %v", err)
	}
	got, err := parseWinner(candOut)
	if err != nil {
		return fmt.Errorf("invalid candidate output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %s, got %s\ncandidate output:\n%s", expected, got, candOut)
	}
	return nil
}

func formatInput(n, m, k int, cells [][2]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for _, c := range cells {
		fmt.Fprintf(&sb, "%d %d\n", c[0], c[1])
	}
	return sb.String()
}

func manualTests() []testCase {
	return []testCase{
		{name: "immediate_row_win", input: formatInput(3, 3, 1, [][2]int{{1, 3}})},
		{name: "immediate_col_win", input: formatInput(4, 4, 1, [][2]int{{4, 1}})},
		{name: "no_special", input: formatInput(2, 2, 0, nil)},
		{name: "sample_like", input: formatInput(3, 3, 1, [][2]int{{2, 2}})},
	}
}

func randomTest(name string, rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN-1) + 1
	m := rng.Intn(maxN-1) + 1
	if n == 1 && m == 1 {
		n = 2
	}
	maxCells := n*m - 1
	if maxCells < 0 {
		maxCells = 0
	}
	k := rng.Intn(min(maxCells, 10) + 1)
	cells := make([][2]int, 0, k)
	used := make(map[[2]int]bool)
	for len(cells) < k {
		x := rng.Intn(n) + 1
		y := rng.Intn(m) + 1
		if x == 1 && y == 1 {
			continue
		}
		key := [2]int{x, y}
		if used[key] {
			continue
		}
		used[key] = true
		cells = append(cells, key)
	}
	return testCase{name: name, input: formatInput(n, m, k, cells)}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomTest(fmt.Sprintf("deterministic_%d", idx+1), rng, 10))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomTest(fmt.Sprintf("random_%d", len(tests)+1), rng, 200))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
