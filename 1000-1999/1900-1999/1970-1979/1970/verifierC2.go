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
	outPath := filepath.Join(verifierDir, "ref1970C2.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1970C2.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
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

func parseOutput(out string) (string, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty output")
	}
	val := fields[0]
	if val != "Ron" && val != "Hermione" {
		return "", fmt.Errorf("invalid output: %s", val)
	}
	return val, nil
}

func verifyCase(candidate, reference string, tc testCase) error {
	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	expected, err := parseOutput(refOut)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate error: %v", err)
	}
	got, err := parseOutput(candOut)
	if err != nil {
		return fmt.Errorf("invalid candidate output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %s, got %s\ncandidate output:\n%s", expected, got, candOut)
	}
	return nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "line3_start1", input: "3 1\n1 2\n2 3\n1\n"},
		{name: "line3_start2", input: "3 1\n1 2\n2 3\n2\n"},
		{name: "star4_start_center", input: "4 1\n1 2\n1 3\n1 4\n1\n"},
	}
}

type edge struct {
	u, v int
}

func randomTree(n int, rng *rand.Rand) []edge {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		edges = append(edges, edge{u, i})
	}
	return edges
}

func randomTest(name string, rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN-1) + 2
	edges := randomTree(n, rng)
	start := rng.Intn(n) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d 1\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	fmt.Fprintf(&sb, "%d\n", start)
	return testCase{name: name, input: sb.String()}
}

func generateTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomTest(fmt.Sprintf("deterministic_%d", idx+1), rng, 15))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 100 {
		maxN := 200
		if len(tests)%10 == 0 {
			maxN = 2000
		}
		tests = append(tests, randomTest(fmt.Sprintf("random_%d", len(tests)+1), rng, maxN))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
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
