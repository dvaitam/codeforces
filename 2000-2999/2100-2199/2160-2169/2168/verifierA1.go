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

const (
	maxN        = 10000
	maxValue    = 26
	maxStrLen   = 100000
	minStrLen   = 1
	firstPhase  = "first"
	secondPhase = "second"
)

type testCase struct {
	name string
	arr  []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2168A1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA1")
	cmd := exec.Command("go", "build", "-o", outPath, "2168A1.go")
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
	return stdout.String(), nil
}

func buildFirstInput(tc testCase) string {
	var arrLine strings.Builder
	for i, v := range tc.arr {
		if i > 0 {
			arrLine.WriteByte(' ')
		}
		arrLine.WriteString(strconv.Itoa(v))
	}
	return fmt.Sprintf("%s\n%d\n%s\n", firstPhase, len(tc.arr), arrLine.String())
}

func firstNonEmptyLine(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return strings.TrimSpace(output)
}

func validateEncodedString(s string) error {
	if len(s) < minStrLen || len(s) > maxStrLen {
		return fmt.Errorf("encoded string length %d not in [%d,%d]", len(s), minStrLen, maxStrLen)
	}
	for i, ch := range s {
		if ch < 'a' || ch > 'z' {
			return fmt.Errorf("encoded string contains invalid character %q at position %d", ch, i+1)
		}
	}
	return nil
}

func parseDecodedOutput(output string, expectedLen int) ([]int, error) {
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("decoding output is empty")
	}
	nVal, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, fmt.Errorf("first token %q is not an integer: %v", tokens[0], err)
	}
	if nVal < 0 || nVal > maxN {
		return nil, fmt.Errorf("decoded n=%d out of allowed range [0,%d]", nVal, maxN)
	}
	if nVal != expectedLen {
		return nil, fmt.Errorf("decoded length %d does not match expected %d", nVal, expectedLen)
	}
	if len(tokens)-1 != nVal {
		return nil, fmt.Errorf("expected %d array values, got %d tokens", nVal, len(tokens)-1)
	}
	arr := make([]int, nVal)
	for i := 0; i < nVal; i++ {
		val, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return nil, fmt.Errorf("token %q at position %d is not an integer: %v", tokens[i+1], i+2, err)
		}
		if val < 1 || val > maxValue {
			return nil, fmt.Errorf("decoded value %d at index %d out of range [1,%d]", val, i, maxValue)
		}
		arr[i] = val
	}
	return arr, nil
}

func verifySingleTest(bin string, tc testCase) error {
	firstInput := buildFirstInput(tc)
	firstOutput, err := runBinary(bin, firstInput)
	if err != nil {
		return fmt.Errorf("encoding run failed: %v", err)
	}
	encoded := firstNonEmptyLine(firstOutput)
	if encoded == "" {
		return fmt.Errorf("encoding run produced empty string")
	}
	if err := validateEncodedString(encoded); err != nil {
		return err
	}

	secondInput := fmt.Sprintf("%s\n%s\n", secondPhase, encoded)
	secondOutput, err := runBinary(bin, secondInput)
	if err != nil {
		return fmt.Errorf("decoding run failed: %v", err)
	}
	decoded, err := parseDecodedOutput(secondOutput, len(tc.arr))
	if err != nil {
		return err
	}
	for i, v := range decoded {
		if v != tc.arr[i] {
			return fmt.Errorf("decoded value mismatch at index %d: expected %d got %d", i, tc.arr[i], v)
		}
	}
	return nil
}

func verifySolution(bin string, tests []testCase) error {
	for idx, tc := range tests {
		if err := verifySingleTest(bin, tc); err != nil {
			return fmt.Errorf("test %d (%s) failed: %w", idx+1, tc.name, err)
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "n1_min", arr: []int{1}},
		{name: "n1_max", arr: []int{maxValue}},
		{name: "all_same_small", arr: []int{5, 5, 5, 5}},
		{name: "ascending_10", arr: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{name: "descending_26", arr: []int{26, 25, 24, 23, 22, 21}},
		{name: "alternating_edges", arr: []int{1, 26, 1, 26, 1, 26, 1, 26}},
		{name: "palindrome", arr: []int{3, 14, 15, 9, 26, 9, 15, 14, 3}},
	}
}

func randomArray(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(maxValue) + 1
	}
	return arr
}

func generateTests() []testCase {
	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i++ {
		n := rng.Intn(50) + 1
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_small_%d", i+1),
			arr:  randomArray(n, rng),
		})
	}

	for i := 0; i < 5; i++ {
		n := rng.Intn(1000) + 100
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_medium_%d", i+1),
			arr:  randomArray(n, rng),
		})
	}

	tests = append(tests,
		testCase{name: "max_length_random", arr: randomArray(maxN, rng)},
		testCase{name: "max_length_pattern", arr: func() []int {
			arr := make([]int, maxN-1)
			for i := range arr {
				arr[i] = (i%maxValue + 1)
			}
			return arr
		}()},
		testCase{name: "half_length_pattern", arr: func() []int {
			n := 5432
			arr := make([]int, n)
			for i := 0; i < n; i++ {
				if i%3 == 0 {
					arr[i] = 1
				} else if i%3 == 1 {
					arr[i] = 13
				} else {
					arr[i] = 26
				}
			}
			return arr
		}()},
	)

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	tests := generateTests()

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	if err := verifySolution(oracle, tests); err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\n", err)
		os.Exit(1)
	}

	if err := verifySolution(target, tests); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
