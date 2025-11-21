package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	maxValue = int64(1_000_000_000_000_000_000)
	maxB     = int64(1_000_000_000)
)

type testCase struct {
	n int
	b []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := buildTestCases(rng)
	input := buildInput(cases)

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	if err := checkProgram(refBin, input, cases); err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	if err := checkProgram(candidate, input, cases); err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2077A.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2077A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func checkProgram(target, input string, cases []testCase) error {
	output, err := runProgram(target, input)
	if err != nil {
		return err
	}
	return validateOutput(output, cases)
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTestCases(rng *rand.Rand) []testCase {
	cases := deterministicCases()
	totalN := 0
	for _, tc := range cases {
		totalN += tc.n
	}

	large := randomCase(rng, 50000)
	cases = append(cases, large)
	totalN += large.n

	for len(cases) < 120 && totalN < 120000 {
		n := rng.Intn(1500) + 1
		if totalN+n > 120000 {
			n = 120000 - totalN
		}
		if n <= 0 {
			break
		}
		cases = append(cases, randomCase(rng, n))
		totalN += n
	}
	return cases
}

func deterministicCases() []testCase {
	return []testCase{
		{n: 1, b: []int64{1, 2}},
		{n: 1, b: []int64{123456789, 987654321}},
		{n: 2, b: []int64{4, 7, 1, 10}},
		{n: 3, b: []int64{5, 2, 9, 14, 1, 7}},
		{n: 4, b: []int64{11, 22, 33, 44, 55, 66, 77, 88}},
	}
}

func randomCase(rng *rand.Rand, n int) testCase {
	return testCase{
		n: n,
		b: distinctBlock(rng, 2*n),
	}
}

func distinctBlock(rng *rand.Rand, length int) []int64 {
	if length <= 0 {
		return nil
	}
	maxStart := maxB - int64(length) + 1
	if maxStart < 1 {
		maxStart = 1
	}
	base := rng.Int63n(maxStart) + 1
	arr := make([]int64, length)
	for i := 0; i < length; i++ {
		arr[i] = base + int64(i)
	}
	rng.Shuffle(length, func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func validateOutput(output string, cases []testCase) error {
	tokens := strings.Fields(output)
	idx := 0
	for caseIdx, tc := range cases {
		expected := 2*tc.n + 1
		if idx+expected > len(tokens) {
			return fmt.Errorf("test %d: expected %d numbers, got %d", caseIdx+1, expected, len(tokens)-idx)
		}
		arr := make([]int64, expected)
		for i := 0; i < expected; i++ {
			val, err := strconv.ParseInt(tokens[idx+i], 10, 64)
			if err != nil {
				return fmt.Errorf("test %d: invalid integer %q (%v)", caseIdx+1, tokens[idx+i], err)
			}
			arr[i] = val
		}
		idx += expected
		if err := validateCase(tc, arr); err != nil {
			return fmt.Errorf("test %d: %v", caseIdx+1, err)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("extra output detected (%d unread tokens)", len(tokens)-idx)
	}
	return nil
}

func validateCase(tc testCase, arr []int64) error {
	if len(arr) != 2*tc.n+1 {
		return fmt.Errorf("expected %d numbers, got %d", 2*tc.n+1, len(arr))
	}
	seen := make(map[int64]struct{}, len(arr))
	for i, v := range arr {
		if v < 1 || v > maxValue {
			return fmt.Errorf("value #%d = %d outside [1, %d]", i+1, v, maxValue)
		}
		if _, ok := seen[v]; ok {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = struct{}{}
	}
	for _, v := range tc.b {
		if _, ok := seen[v]; !ok {
			return fmt.Errorf("missing original number %d", v)
		}
	}
	if len(seen) != len(tc.b)+1 {
		return fmt.Errorf("expected %d distinct values, got %d", len(tc.b)+1, len(seen))
	}

	sum := big.NewInt(0)
	tmp := big.NewInt(0)
	for i := 1; i < len(arr); i++ {
		tmp.SetInt64(arr[i])
		if i%2 == 1 {
			sum.Add(sum, tmp)
		} else {
			sum.Sub(sum, tmp)
		}
	}
	first := big.NewInt(0).SetInt64(arr[0])
	if sum.Cmp(first) != 0 {
		return fmt.Errorf("fails alternating equality: a1=%s, computed=%s", first.String(), sum.String())
	}
	return nil
}
