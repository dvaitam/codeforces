package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const tol = 1e-6

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		want, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}
		if !closeEnough(want, got) {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %.10f, got %.10f\nInput:\n%s\nCandidate output:\n%s\n", idx+1, want, got, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1765F.go",
		filepath.Join("1000-1999", "1700-1799", "1760-1769", "1765", "1765F.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1765F.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1765F_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q", fields[0])
	}
	return val, nil
}

func closeEnough(a, b float64) bool {
	diff := math.Abs(a - b)
	if diff <= tol {
		return true
	}
	if math.Abs(b) > 1 {
		return diff/math.Abs(b) <= tol
	}
	return false
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([]contract{{0, 1, 1}}, 1),
		buildTest([]contract{{0, 10, 100}, {100, 10, 200}}, 5),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(6) + 1
		tests = append(tests, randomTest(rng, n))
	}
	tests = append(tests, randomTest(rng, 50))
	return tests
}

type contract struct {
	x int
	w int
	c int
}

func buildTest(cs []contract, k int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", len(cs), k))
	for _, c := range cs {
		b.WriteString(fmt.Sprintf("%d %d %d\n", c.x, c.w, c.c))
	}
	return testCase{input: b.String()}
}

func randomTest(rng *rand.Rand, n int) testCase {
	k := rng.Intn(100) + 1
	cs := make([]contract, n)
	for i := 0; i < n; i++ {
		cs[i] = contract{
			x: rng.Intn(101),
			w: rng.Intn(1000) + 1,
			c: rng.Intn(1000) + 1,
		}
	}
	return buildTest(cs, k)
}
