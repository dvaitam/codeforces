package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceD = "2000-2999/2100-2199/2170-2179/2172/2172D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceD)
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
		refVals, err := parseModInts(refOut, input.students)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, input.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input.input)
			os.Exit(1)
		}
		userVals, err := parseModInts(userOut, input.students)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input.input, userOut)
			os.Exit(1)
		}

		if err := compareModAnswers(refVals, userVals); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", idx+1, err, input.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
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
	tmp, err := os.CreateTemp("", "2172D-ref-*")
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

type testInput struct {
	input    string
	students int
}

func parseModInts(out string, expected int) ([]*big.Int, error) {
	reader := strings.NewReader(out)
	res := make([]*big.Int, expected)
	for i := 0; i < expected; i++ {
		var s string
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return nil, fmt.Errorf("expected %d numbers, got %d (%v)", expected, i, err)
		}
		val, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		res[i] = val.Mod(val, big.NewInt(998244353))
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func compareModAnswers(expected, actual []*big.Int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d answers, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i].Cmp(actual[i]) != 0 {
			return fmt.Errorf("answer %d mismatch: expected %s, got %s", i+1, expected[i].String(), actual[i].String())
		}
	}
	return nil
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTest())
	tests = append(tests, simpleOwnershipTest())
	rng := rand.New(rand.NewSource(2172))
	tests = append(tests, randomTest(rng, 3, 2))
	tests = append(tests, randomTest(rng, 10, 5))
	tests = append(tests, randomTest(rng, 50, 20))
	tests = append(tests, randomTest(rng, 200, 100))
	tests = append(tests, randomTest(rng, 600, 400))
	return tests
}

func sampleTest() testInput {
	input := strings.TrimSpace(`5 2
1 2 3 4 5
0 0 1 2 1
`) + "\n"
	return testInput{input: input, students: 2}
}

func simpleOwnershipTest() testInput {
	n, m := 6, 3
	a := []int64{1, 2, 4, 8, 16, 32}
	b := []int{0, 1, 2, 3, 0, 0}
	return buildCase(n, m, a, b)
}

func randomTest(rng *rand.Rand, n, m int) testInput {
	if m >= n {
		m = n - 1
	}
	a := make([]int64, n)
	a[0] = int64(rng.Intn(10) + 1)
	for i := 1; i < n; i++ {
		a[i] = a[i-1] + int64(rng.Intn(5)+1)
	}
	b := make([]int, n)
	count := make([]int, m+1)
	for i := 0; i < n; i++ {
		who := rng.Intn(m + 1)
		b[i] = who
		count[who]++
	}
	for j := 0; j <= m; j++ {
		if count[j] == 0 {
			idx := rng.Intn(n)
			b[idx] = j
			count[j]++
		}
	}
	return buildCase(n, m, a, b)
}

func buildCase(n, m int, a []int64, b []int) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	return testInput{input: sb.String(), students: m}
}
