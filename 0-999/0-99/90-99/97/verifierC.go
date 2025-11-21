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
	"sort"
	"strings"
	"time"
)

const tolerance = 1e-6

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
	ref := filepath.Join(verifierDir, "ref97C.bin")
	cmd := exec.Command("go", "build", "-o", ref, "97C.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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

func parseFloatOutput(out string) (float64, error) {
	reader := strings.NewReader(out)
	var val float64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("failed to read float: %v\n%s", err, out)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("output is not a finite number: %v", val)
	}
	return val, nil
}

func almostEqual(got, expected float64) bool {
	limit := tolerance * math.Max(1.0, math.Abs(expected))
	return math.Abs(got-expected) <= limit+1e-12
}

func verifyCase(candidate, reference string, tc testCase) error {
	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expected, err := parseFloatOutput(refOut)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return err
	}
	got, err := parseFloatOutput(candOut)
	if err != nil {
		return err
	}
	if !almostEqual(got, expected) {
		return fmt.Errorf("expected %.12f, got %.12f", expected, got)
	}
	return nil
}

func formatInput(n int, probs []float64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i <= n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%.6f", probs[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func manualTest(name string, probs []float64) testCase {
	return testCase{
		name:  name,
		input: formatInput(len(probs)-1, probs),
	}
}

func randomProbabilities(rng *rand.Rand, n int) []float64 {
	vals := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		vals[i] = rng.Float64()
	}
	sort.Float64s(vals)
	for i := 0; i <= n; i++ {
		vals[i] = math.Round(vals[i]*1e6) / 1e6
		if i > 0 && vals[i] < vals[i-1] {
			vals[i] = vals[i-1]
		}
	}
	if vals[0] < 0 {
		vals[0] = 0
	}
	if vals[n] > 1 {
		vals[n] = 1
	}
	return vals
}

func randomCase(name string, seed int64, n int) testCase {
	rng := rand.New(rand.NewSource(seed))
	return testCase{
		name:  name,
		input: formatInput(n, randomProbabilities(rng, n)),
	}
}

func generateTests() []testCase {
	tests := []testCase{
		manualTest("all_zero", []float64{0, 0, 0, 0}),
		manualTest("all_one", []float64{1, 1, 1, 1}),
		manualTest("increasing", []float64{0.0, 0.1, 0.2, 0.4, 0.6}),
		manualTest("plateau", []float64{0.2, 0.2, 0.2, 0.2, 0.9, 0.9}),
		manualTest("sharp_jump", []float64{0, 0.001, 0.002, 0.9, 0.95, 0.99, 1.0}),
		manualTest("linear100", buildLinear(100)),
	}
	deterministicSeeds := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, seed := range deterministicSeeds {
		n := 3 + int(seed)%98
		tests = append(tests, randomCase(fmt.Sprintf("deterministic_%d", seed), seed, n))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		n := rng.Intn(98) + 3
		seed := rng.Int63()
		tests = append(tests, randomCase(fmt.Sprintf("random_%d", seed), seed, n))
	}
	return tests
}

func buildLinear(n int) []float64 {
	vals := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		vals[i] = math.Round(float64(i)/float64(n)*1e6) / 1e6
	}
	return vals
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
			name := tc.name
			if name == "" {
				name = fmt.Sprintf("case_%d", i+1)
			}
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", i+1, name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
