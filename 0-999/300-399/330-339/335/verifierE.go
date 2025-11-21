package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "0-999/300-399/330-339/335/335E.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		local := filepath.Join(dir, "335E.go")
		if _, err := os.Stat(local); err == nil {
			referenceSolutionPath = local
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
	input string
}

func formatAlice(n, h int) string {
	return fmt.Sprintf("Alice\n%d %d\n", n, h)
}

func formatBob(counter, h int) string {
	return fmt.Sprintf("Bob\n%d %d\n", counter, h)
}

func generateTests() []testCase {
	tests := []testCase{
		{input: formatAlice(3, 1)},
		{input: formatBob(2, 30)},
		{input: formatAlice(2572, 10)},
		{input: formatAlice(2, 0)},
		{input: formatAlice(30000, 30)},
		{input: formatBob(2, 0)},
		{input: formatBob(30000, 30)},
	}
	rng := rand.New(rand.NewSource(20250305))
	for i := 0; i < 150; i++ {
		n := rng.Intn(30000-2+1) + 2
		h := rng.Intn(31)
		tests = append(tests, testCase{input: formatAlice(n, h)})
	}
	for i := 0; i < 150; i++ {
		counter := rng.Intn(30000-2+1) + 2
		h := rng.Intn(31)
		tests = append(tests, testCase{input: formatBob(counter, h)})
	}
	for i := 0; i < 50; i++ {
		n := 2 + i%5
		h := i % 31
		tests = append(tests, testCase{input: formatAlice(n, h)})
		tests = append(tests, testCase{input: formatBob(n, h)})
	}
	return tests
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
		return "", nil, fmt.Errorf("reference solution path is empty")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "335E-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_335E")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference solution: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseFloatOutput(out string) (float64, error) {
	reader := strings.NewReader(strings.TrimSpace(out))
	var val float64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("failed to parse float from output %q: %v", out, err)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("output is not a finite number: %v", val)
	}
	return val, nil
}

func almostEqual(a, b float64) bool {
	diff := math.Abs(a - b)
	const absTol = 1e-7
	const relTol = 1e-7
	if diff <= absTol {
		return true
	}
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return diff <= absTol
	}
	return diff <= relTol*scale
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, refOut)
			os.Exit(1)
		}
		refVal, err := parseFloatOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, refOut)
			os.Exit(1)
		}
		out, runErr := runProgram(bin, tc.input)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, tc.input, out)
			os.Exit(1)
		}
		gotVal, err := parseFloatOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
		if !almostEqual(refVal, gotVal) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected: %.10f\ngot: %.10f\nraw output:\n%s\n", i+1, tc.input, refVal, gotVal, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
