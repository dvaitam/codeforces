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

const referenceSource = "2000-2999/2100-2199/2100-2109/2101/2101A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	candidatePath, err := candidatePathFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	candidate, candCleanup, err := buildCandidate(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	refPath := referencePath()
	refBin, refCleanup, err := buildReferenceBinary(refPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := buildTests()
	for idx, tc := range tests {
		ns, err := parseInputNs(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing generated input %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		refGrids, err := parseGrids(refOut, ns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\nraw:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}
		refScores := make([]int64, len(ns))
		for i, g := range refGrids {
			score, err := computeSum(ns[i], g)
			if err != nil {
				fmt.Fprintf(os.Stderr, "reference grid invalid on test %d (%s) case %d: %v\n", idx+1, tc.name, i+1, err)
				os.Exit(1)
			}
			refScores[i] = score
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		candGrids, err := parseGrids(candOut, ns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\nraw:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		for i, g := range candGrids {
			score, err := computeSum(ns[i], g)
			if err != nil {
				fmt.Fprintf(os.Stderr, "candidate grid invalid on test %d (%s) case %d: %v\n", idx+1, tc.name, i+1, err)
				fmt.Fprintln(os.Stderr, previewInput(tc.input))
				os.Exit(1)
			}
			if score != refScores[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected score %d got %d\n", idx+1, tc.name, i+1, refScores[i], score)
				fmt.Fprintln(os.Stderr, previewInput(tc.input))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("usage: go run verifierA.go /path/to/binary-or-source")
	}
	return os.Args[1], nil
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), "2101A.go")
	}
	return referenceSource
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2101A-ref")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2101a")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func buildCandidate(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "verifier2101A-cand")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", bin, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseInputNs(in string) ([]int, error) {
	tokens := strings.Fields(in)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	if t < 1 || t > 100 {
		return nil, fmt.Errorf("t out of bounds: %d", t)
	}
	if len(tokens) != t+1 {
		return nil, fmt.Errorf("expected %d n values, got %d", t, len(tokens)-1)
	}
	ns := make([]int, t)
	sumN := 0
	for i := 0; i < t; i++ {
		n, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid n at case %d: %v", i+1, err)
		}
		if n < 1 || n > 500 {
			return nil, fmt.Errorf("n out of bounds at case %d: %d", i+1, n)
		}
		sumN += n
		if sumN > 1000 {
			return nil, fmt.Errorf("sum of n exceeds 1000")
		}
		ns[i] = n
	}
	return ns, nil
}

func parseGrids(out string, ns []int) ([][]int, error) {
	tokens := strings.Fields(out)
	total := 0
	for _, n := range ns {
		total += n * n
	}
	if len(tokens) != total {
		return nil, fmt.Errorf("expected %d grid integers, got %d", total, len(tokens))
	}
	res := make([][]int, len(ns))
	idx := 0
	for i, n := range ns {
		sz := n * n
		grid := make([]int, sz)
		for j := 0; j < sz; j++ {
			val, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", tokens[idx])
			}
			grid[j] = val
			idx++
		}
		res[i] = grid
	}
	return res, nil
}

func computeSum(n int, grid []int) (int64, error) {
	if len(grid) != n*n {
		return 0, fmt.Errorf("expected %d values, got %d", n*n, len(grid))
	}
	n2 := n * n
	row := make([]int, n2)
	col := make([]int, n2)
	for i := 0; i < n2; i++ {
		row[i] = -1
		col[i] = -1
	}
	for idx, val := range grid {
		if val < 0 || val >= n2 {
			return 0, fmt.Errorf("value out of range: %d", val)
		}
		if row[val] != -1 {
			return 0, fmt.Errorf("duplicate value %d", val)
		}
		row[val] = idx / n
		col[val] = idx % n
	}
	minR, maxR := row[0], row[0]
	minC, maxC := col[0], col[0]
	sum := int64(0)
	for k := 1; k <= n2; k++ {
		if row[k-1] < minR {
			minR = row[k-1]
		}
		if row[k-1] > maxR {
			maxR = row[k-1]
		}
		if col[k-1] < minC {
			minC = col[k-1]
		}
		if col[k-1] > maxC {
			maxC = col[k-1]
		}
		topChoices := int64(minR + 1)
		bottomChoices := int64(n - maxR)
		leftChoices := int64(minC + 1)
		rightChoices := int64(n - maxC)
		sum += topChoices * bottomChoices * leftChoices * rightChoices
	}
	return sum, nil
}

func buildTests() []testCase {
	var tests []testCase

	tests = append(tests, testCase{
		name:  "small-sample",
		input: "2\n2\n3\n",
	})

	tests = append(tests, testCase{
		name:  "minimal",
		input: "1\n1\n",
	})

	tests = append(tests, testCase{
		name:  "medium",
		input: "2\n4\n5\n",
	})

	tests = append(tests, testCase{
		name:  "large-single",
		input: "1\n200\n",
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tests = append(tests, randomTest(rng, i, 80))
	}

	return tests
}

func randomTest(rng *rand.Rand, idx, maxN int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	sumN := 0
	for i := 0; i < t; i++ {
		maxAllowed := 1000 - sumN - (t - i - 1)
		if maxAllowed < 1 {
			maxAllowed = 1
		}
		n := rng.Intn(min(maxN, maxAllowed)) + 1
		if n > 500 {
			n = 500
		}
		sumN += n
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return testCase{
		name:  fmt.Sprintf("random-%d", idx+1),
		input: sb.String(),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func previewInput(in string) string {
	const limit = 500
	if len(in) <= limit {
		return in
	}
	return in[:limit] + "..."
}
