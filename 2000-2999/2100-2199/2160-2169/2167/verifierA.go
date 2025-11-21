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

type testCase struct {
	name  string
	input string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2167A-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, "2167A.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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
	return strings.TrimSpace(stdout.String()), nil
}

func makeInput(cases [][4]int) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", c[0], c[1], c[2], c[3]))
	}
	return sb.String()
}

func sampleTests() []testCase {
	cases := [][4]int{
		{1, 2, 3, 4},
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{5, 5, 5, 1},
		{4, 10, 5, 9},
	}
	return []testCase{
		{
			name:  "samples",
			input: makeInput(cases),
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 20)
	for i := 0; i < 20; i++ {
		count := rng.Intn(500) + 1
		cases := make([][4]int, count)
		for j := 0; j < count; j++ {
			for k := 0; k < 4; k++ {
				cases[j][k] = rng.Intn(10) + 1
			}
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: makeInput(cases),
		})
	}
	return tests
}

func exhaustiveTest() testCase {
	cases := make([][4]int, 0, 10000)
	for a := 1; a <= 10; a++ {
		for b := 1; b <= 10; b++ {
			for c := 1; c <= 10; c++ {
				for d := 1; d <= 10; d++ {
					cases = append(cases, [4]int{a, b, c, d})
				}
			}
		}
	}
	return testCase{
		name:  "exhaustive_all",
		input: makeInput(cases),
	}
}

func normalizeAnswers(output string) []string {
	fields := strings.Fields(output)
	res := make([]string, len(fields))
	for i, f := range fields {
		res[i] = strings.ToUpper(f)
	}
	return res
}

func compareOutputs(expected, actual string) error {
	exp := normalizeAnswers(expected)
	act := normalizeAnswers(actual)
	if len(exp) != len(act) {
		return fmt.Errorf("token count mismatch: expected %d got %d", len(exp), len(act))
	}
	for i := range exp {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at token %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(sampleTests(), randomTests()...)
	tests = append(tests, exhaustiveTest())

	for idx, tc := range tests {
		expected, err := runBinary(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		actual, err := runBinary(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := compareOutputs(expected, actual); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.name, err, tc.input, expected, actual)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
