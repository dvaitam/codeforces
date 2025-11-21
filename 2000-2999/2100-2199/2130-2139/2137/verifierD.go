package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2130-2139/2137/2137D.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-2137D-*")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

type testCase struct {
	n int
	b []int
}

func validateCandidate(tc testCase, candTokens []int) error {
	// if candTokens length 1 and token == -1 => only acceptable if no solution exists
	if len(candTokens) == 1 && candTokens[0] == -1 {
		// handled elsewhere by comparing with reference feasibility
		return nil
	}
	if len(candTokens) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(candTokens))
	}
	for _, v := range candTokens {
		if v < 1 || v > tc.n {
			return fmt.Errorf("value %d out of range [1,%d]", v, tc.n)
		}
	}
	freq := make(map[int]int)
	for _, v := range candTokens {
		freq[v]++
	}
	for i, v := range candTokens {
		if freq[v] != tc.b[i] {
			return fmt.Errorf("position %d expects freq %d for value %d, got %d", i, tc.b[i], v, freq[v])
		}
	}
	return nil
}

func consumeOutput(tc testCase, r *strings.Reader) ([]int, error) {
	// peek next token; if -1 accept single token
	var tokens []int
	var first int
	if _, err := fmt.Fscan(r, &first); err != nil {
		return nil, fmt.Errorf("cannot read answer: %v", err)
	}
	if first == -1 {
		return []int{-1}, nil
	}
	tokens = append(tokens, first)
	for len(tokens) < tc.n {
		var x int
		if _, err := fmt.Fscan(r, &x); err != nil {
			return nil, fmt.Errorf("not enough numbers: %v", err)
		}
		tokens = append(tokens, x)
	}
	return tokens, nil
}

func fixedTests() []string {
	return []string{
		"3\n4\n1 2 3 4\n6\n1 2 2 3 3 3\n6\n6 6 6 6 6 6\n",
		"2\n1\n1\n5\n1 1 1 1 1\n",
	}
}

func generateValidCase(rng *rand.Rand, n int) testCase {
	// Create groups with size v, repeated v times
	var b []int
	labelSizes := []int{}
	for len(b) < n {
		size := rng.Intn(6) + 1 // 1..6
		if len(b)+size > n {
			size = n - len(b)
		}
		labelSizes = append(labelSizes, size)
		for i := 0; i < size; i++ {
			b = append(b, size)
		}
	}
	// shuffle b to avoid organized pattern
	for i := range b {
		j := rng.Intn(i + 1)
		b[i], b[j] = b[j], b[i]
	}
	return testCase{n: n, b: b}
}

func generateInvalidCase(rng *rand.Rand, n int) testCase {
	// start with valid then perturb count for one value to break divisibility
	tc := generateValidCase(rng, n)
	if n == 1 {
		tc.b[0] = rng.Intn(2) + 2 // set to 2 or 3 => impossible
		return tc
	}
	// pick value v and change a count so occurrences not divisible by v
	idx := rng.Intn(n)
	tc.b[idx] = tc.b[idx] + 1
	return tc
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomTests(rng *rand.Rand) string {
	t := rng.Intn(20) + 1
	var cases []testCase
	totalN := 0
	for i := 0; i < t; i++ {
		if totalN >= 150000 {
			break
		}
		n := rng.Intn(50) + 1
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		totalN += n
		var tc testCase
		if rng.Intn(3) == 0 {
			tc = generateInvalidCase(rng, n)
		} else {
			tc = generateValidCase(rng, n)
		}
		cases = append(cases, tc)
	}
	return buildInput(cases)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refPath, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTests(rng))
	}

	for idx, input := range tests {
		var t int
		if _, err := fmt.Fscan(strings.NewReader(input), &t); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse t for case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		// parse input test cases for validation
		inReader := strings.NewReader(input)
		fmt.Fscan(inReader, new(int)) // skip t
		testCases := make([]testCase, 0, t)
		for i := 0; i < t; i++ {
			var n int
			fmt.Fscan(inReader, &n)
			tc := testCase{n: n, b: make([]int, n)}
			for j := 0; j < n; j++ {
				fmt.Fscan(inReader, &tc.b[j])
			}
			testCases = append(testCases, tc)
		}

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		candOut, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}

		refReader := strings.NewReader(refOut)
		candReader := strings.NewReader(candOut)
		for caseIdx, tc := range testCases {
			refAns, err := consumeOutput(tc, refReader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d case %d: %v\noutput:\n%s", idx+1, caseIdx+1, err, refOut)
				os.Exit(1)
			}
			candAns, err := consumeOutput(tc, candReader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d case %d: %v\n", idx+1, caseIdx+1, err)
				os.Exit(1)
			}
			refFeasible := !(len(refAns) == 1 && refAns[0] == -1)
			if !refFeasible {
				if !(len(candAns) == 1 && candAns[0] == -1) {
					fmt.Fprintf(os.Stderr, "candidate should output -1 on impossible case (test %d case %d)\ninput:\n%s", idx+1, caseIdx+1, input)
					os.Exit(1)
				}
				continue
			}
			if len(candAns) == 1 && candAns[0] == -1 {
				fmt.Fprintf(os.Stderr, "candidate reported -1 but solution exists (test %d case %d)\ninput:\n%s", idx+1, caseIdx+1, input)
				os.Exit(1)
			}
			if err := validateCandidate(tc, candAns); err != nil {
				fmt.Fprintf(os.Stderr, "invalid candidate on test %d case %d: %v\ninput:\n%s", idx+1, caseIdx+1, err, input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed")
}
