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

const refSource2038A = "./2038A.go"

type testCase struct {
	n int
	k int64
	a []int64
	b []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2038A)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		refOut, err := runProgram(refBin, input.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input.input)
			os.Exit(1)
		}
		refVals, err := parseInt64s(refOut, input.caseData.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, input.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input.input)
			os.Exit(1)
		}
		userVals, err := parseInt64s(userOut, input.caseData.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input.input, userOut)
			os.Exit(1)
		}

		if err := validateSolution(input.caseData, userVals, refVals); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", idx+1, err, input.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

type testInput struct {
	input    string
	caseData testCase
}

func validateSolution(tc testCase, sol []int64, optimal []int64) error {
	if len(sol) != tc.n {
		return fmt.Errorf("expected %d values, got %d", tc.n, len(sol))
	}

	total := int64(0)
	for i := 0; i < tc.n; i++ {
		if sol[i] < 0 || sol[i] > tc.k {
			return fmt.Errorf("invalid work amount c[%d]=%d", i+1, sol[i])
		}
		total += sol[i]
	}

	limit := make([]int64, tc.n)
	totalLimit := int64(0)
	for i := 0; i < tc.n; i++ {
		if tc.b[i] == 0 {
			limit[i] = tc.k
		} else {
			limit[i] = tc.a[i] / tc.b[i]
		}
		totalLimit += limit[i]
	}
	if totalLimit < tc.k {
		for i := 0; i < tc.n; i++ {
			if sol[i] != 0 {
				return fmt.Errorf("project impossible, but c[%d]=%d", i+1, sol[i])
			}
		}
		return nil
	}

	required := tc.k
	for i := 0; i < tc.n; i++ {
		need := required
		for j := i + 1; j < tc.n; j++ {
			need -= limit[j]
			if need <= 0 {
				need = 0
				break
			}
		}
		if need < 0 {
			need = 0
		}
		if need > limit[i] {
			need = limit[i]
		}
		if sol[i] != need {
			return fmt.Errorf("engineer %d should announce %d, got %d", i+1, need, sol[i])
		}
		required -= sol[i]
	}
	return nil
}

func parseInt64s(out string, n int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("expected %d numbers, got %d (%v)", n, i, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2038A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(source))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTest())
	tests = append(tests, impossibleCase())
	tests = append(tests, zeroBenefitCase())
	rng := rand.New(rand.NewSource(2038))
	tests = append(tests, randomTest(rng, 5, 20, 50))
	tests = append(tests, randomTest(rng, 20, 200, 1000))
	tests = append(tests, randomTest(rng, 1000, 1000000, 1000000))
	return tests
}

func sampleTest() testInput {
	tc := testCase{
		n: 3,
		k: 6,
		a: []int64{4, 7, 6},
		b: []int64{1, 2, 3},
	}
	return buildTest(tc)
}

func impossibleCase() testInput {
	tc := testCase{
		n: 3,
		k: 12,
		a: []int64{4, 7, 6},
		b: []int64{1, 2, 3},
	}
	return buildTest(tc)
}

func zeroBenefitCase() testInput {
	tc := testCase{
		n: 3,
		k: 11,
		a: []int64{6, 7, 8},
		b: []int64{1, 2, 3},
	}
	return buildTest(tc)
}

func randomTest(rng *rand.Rand, n int, maxK int64, maxA int64) testInput {
	tc := testCase{
		n: n,
		k: rng.Int63n(maxK) + 1,
		a: make([]int64, n),
		b: make([]int64, n),
	}
	for i := 0; i < n; i++ {
		tc.a[i] = rng.Int63n(maxA) + 1
		tc.b[i] = rng.Int63n(1000) + 1
	}
	return buildTest(tc)
}

func buildTest(tc testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.b[i]))
	}
	sb.WriteByte('\n')
	return testInput{input: sb.String(), caseData: tc}
}
