package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "./2094C.go"

type testCase struct {
	n    int
	grid []string
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2094C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func permutationToGrid(perm []int) []string {
	n := len(perm) / 2
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		builder := strings.Builder{}
		for j := 0; j < n; j++ {
			if j > 0 {
				builder.WriteByte(' ')
			}
			idx := i + j + 2 // p_{i+j} with 1-indexed i,j => index = i+j
			val := perm[idx]
			builder.WriteString(strconv.Itoa(val))
		}
		grid[i] = builder.String()
	}
	return grid
}

func makeTestFromPerm(perm []int) testCase {
	return testCase{n: len(perm) / 2, grid: permutationToGrid(perm)}
}

func deterministicTests() []testCase {
	tests := []testCase{}

	// n = 1
	tests = append(tests, makeTestFromPerm([]int{0, 1, 2}))

	// n = 2, permutation length 4
	tests = append(tests, makeTestFromPerm([]int{0, 2, 1, 3, 4}))

	// n = 3, sample-like permutation
	tests = append(tests, makeTestFromPerm([]int{0, 5, 1, 6, 2, 4, 3}))

	return tests
}

func randomPermutation(rng *rand.Rand, n int) []int {
	size := 2 * n
	perm := make([]int, size+1) // 1-indexed for convenience
	for i := 1; i <= size; i++ {
		perm[i] = i
	}
	for i := size; i >= 2; i-- {
		j := rng.Intn(i) + 1
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm
}

func randomTest(rng *rand.Rand, n int) testCase {
	perm := randomPermutation(rng, n)
	return makeTestFromPerm(perm)
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([][]int, len(tests))
	for i, tc := range tests {
		need := 2 * tc.n
		line := make([]int, need)
		seen := make([]bool, need+1)
		for j := 0; j < need; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("test %d: expected %d numbers, got %d", i+1, need, j)
			}
			val, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, sc.Text())
			}
			if val < 1 || val > need {
				return nil, fmt.Errorf("test %d: value out of range %d", i+1, val)
			}
			if seen[val] {
				return nil, fmt.Errorf("test %d: duplicate value %d", i+1, val)
			}
			seen[val] = true
			line[j] = val
		}
		res[i] = line
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", len(tests))
	}
	return res, nil
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 6; i++ {
		n := rng.Intn(20) + 1
		tests = append(tests, randomTest(rng, n))
	}
	for i := 0; i < 2; i++ {
		n := rng.Intn(60) + 20
		tests = append(tests, randomTest(rng, n))
	}
	// Larger structured cases while keeping sum(n) <= 800
	tests = append(tests, randomTest(rng, 150))
	tests = append(tests, randomTest(rng, 300))

	if totalN(tests) > 800 {
		fmt.Fprintln(os.Stderr, "generated tests exceed constraint on total n")
		os.Exit(1)
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if len(want[i]) != len(got[i]) {
			fmt.Fprintf(os.Stderr, "test %d failed: length mismatch expected %d got %d\n", i+1, len(want[i]), len(got[i]))
			os.Exit(1)
		}
		mismatch := false
		for j := range want[i] {
			if want[i][j] != got[i][j] {
				mismatch = true
				break
			}
		}
		if mismatch {
			fmt.Fprintf(os.Stderr, "test %d failed: outputs differ from reference\nn=%d grid first row=%s\n", i+1, tests[i].n, tests[i].grid[0])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
