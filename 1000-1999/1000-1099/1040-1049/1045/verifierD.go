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

const tol = 1e-4

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		wantVals, err := parseOutputs(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if len(gotVals) != len(wantVals) {
			fmt.Fprintf(os.Stderr, "candidate output count mismatch on test %d\nExpected %d values, got %d\n", idx+1, len(wantVals), len(gotVals))
			os.Exit(1)
		}
		for i := range wantVals {
			if !closeEnough(gotVals[i], wantVals[i]) {
				fmt.Fprintf(os.Stderr, "test %d: value %d differs. expected %.6f, got %.6f\nInput:\n%sCandidate output:\n%s\n", idx+1, i+1, wantVals[i], gotVals[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1045D.go",
		filepath.Join("1000-1999", "1000-1099", "1040-1049", "1045", "1045D.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 1045D.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1045D_%d.bin", time.Now().UnixNano()))
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

func parseOutputs(out string) ([]float64, error) {
	var vals []float64
	lines := strings.Fields(out)
	for _, tok := range lines {
		val, err := strconvParseFloat(tok)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float %q: %v", tok, err)
		}
		vals = append(vals, val)
	}
	return vals, nil
}

func strconvParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func closeEnough(a, b float64) bool {
	diff := math.Abs(a - b)
	if diff <= tol {
		return true
	}
	if math.Abs(b) > 0 {
		if diff/math.Abs(b) <= tol {
			return true
		}
	}
	return false
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(1, []float64{0}, nil, nil, nil),
		buildTest(2, []float64{0.5, 0.5}, [][]int{{1, 2}}, []int{1}, []float64{0.25}),
		buildTest(3, []float64{0.1, 0.2, 0.3}, [][]int{{1, 2}, {2, 3}}, []int{2}, []float64{0.5}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		n := 1 + rng.Intn(200)
		tests = append(tests, randomTest(rng, n, 50))
	}
	tests = append(tests, randomTest(rng, 1000, 1000))
	tests = append(tests, randomTest(rng, 100000, 100000))
	return tests
}

func randomTest(rng *rand.Rand, n, q int) testCase {
	if n < 1 {
		n = 1
	}
	if q < 1 {
		q = 1
	}
	edges := randomTree(rng, n)
	probs := randomProbabilities(rng, n)
	updates := make([]int, q)
	newVals := make([]float64, q)
	for i := 0; i < q; i++ {
		updates[i] = rng.Intn(n) + 1
		newVals[i] = randomProbability(rng)
	}
	return buildTest(n, probs, edges, updates, newVals)
}

func buildTest(n int, probs []float64, edges [][]int, updates []int, newVals []float64) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%.2f", probs[i]))
	}
	b.WriteByte('\n')
	if edges != nil {
		for _, e := range edges {
			b.WriteString(fmt.Sprintf("%d %d\n", e[0]-1, e[1]-1))
		}
	} else {
		for i := 1; i < n; i++ {
			b.WriteString(fmt.Sprintf("%d %d\n", i-1, i))
		}
	}
	if updates == nil {
		updates = []int{1}
		newVals = []float64{probabilitiesClamp(0.5)}
	}
	b.WriteString(fmt.Sprintf("%d\n", len(updates)))
	for i := range updates {
		b.WriteString(fmt.Sprintf("%d %.2f\n", updates[i]-1, newVals[i]))
	}
	return testCase{input: b.String()}
}

func randomTree(rng *rand.Rand, n int) [][]int {
	edges := make([][]int, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := rng.Intn(i-1) + 1
		edges = append(edges, []int{parent, i})
	}
	return edges
}

func randomProbabilities(rng *rand.Rand, n int) []float64 {
	arr := make([]float64, n)
	for i := 0; i < n; i++ {
		arr[i] = randomProbability(rng)
	}
	return arr
}

func randomProbability(rng *rand.Rand) float64 {
	return math.Round(rng.Float64()*100) / 100
}

func probabilitiesClamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
