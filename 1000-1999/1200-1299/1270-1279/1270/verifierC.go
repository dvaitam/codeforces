package main

import (
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

const maxVal int64 = 1e18

type testCase struct {
	input string
	n     int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", out, "1270C.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return out, nil
}

func runProgram(bin string, input string) (string, error) {
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

func parseOutput(out string, t int) ([]int, [][]int64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != t*2 {
		return nil, nil, fmt.Errorf("expected %d lines, got %d", t*2, len(lines))
	}
	counts := make([]int, t)
	appends := make([][]int64, t)
	for i := 0; i < t; i++ {
		lineCnt := strings.TrimSpace(lines[2*i])
		if lineCnt == "" {
			return nil, nil, fmt.Errorf("empty count line at case %d", i+1)
		}
		cnt, err := strconv.Atoi(lineCnt)
		if err != nil || cnt < 0 || cnt > 3 {
			return nil, nil, fmt.Errorf("invalid append count at case %d: %q", i+1, lineCnt)
		}
		counts[i] = cnt
		if cnt == 0 {
			appends[i] = nil
			continue
		}
		values := strings.Fields(lines[2*i+1])
		if len(values) != cnt {
			return nil, nil, fmt.Errorf("case %d: expected %d numbers, got %d", i+1, cnt, len(values))
		}
		app := make([]int64, cnt)
		for j, val := range values {
			num, err := strconv.ParseInt(val, 10, 64)
			if err != nil || num < 0 || num > maxVal {
				return nil, nil, fmt.Errorf("case %d: invalid append value %q", i+1, val)
			}
			app[j] = num
		}
		appends[i] = app
	}
	return counts, appends, nil
}

func deterministicTests() []testCase {
	return []testCase{
		formatTest([]int64{1}),
		formatTest([]int64{0}),
		formatTest([]int64{1, 2, 3}),
		formatTest([]int64{10, 20, 30, 40}),
		formatTest([]int64{1, 1, 1, 1}),
		formatTest([]int64{0, 0, 0}),
	}
}

func formatTest(nums []int64) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(len(nums)))
	sb.WriteByte('\n')
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), n: len(nums)}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	totalN := 0
	for len(tests) < count && totalN < 100000 {
		n := rnd.Intn(100000-len(tests)) + 1
		if totalN+n > 100000 {
			n = 100000 - totalN
		}
		totalN += n
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			x := rnd.Int63n(1_000_000_000 + 1)
			sb.WriteString(strconv.FormatInt(x, 10))
		}
		sb.WriteByte('\n')
		tests = append(tests, testCase{input: sb.String(), n: n})
	}
	return tests
}

func validateSolution(original []int64, appended []int64) bool {
	sum := int64(0)
	xor := int64(0)
	for _, v := range original {
		sum += v
		xor ^= v
	}
	for _, v := range appended {
		sum += v
		xor ^= v
	}
	return sum == xor*2
}

func extractTestData(input string) (int, [][]int64, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, nil, fmt.Errorf("failed to read t: %v", err)
	}
	testData := make([][]int64, t)
	for tc := 0; tc < t; tc++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return 0, nil, fmt.Errorf("case %d: failed to read n: %v", tc+1, err)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
				return 0, nil, fmt.Errorf("case %d: failed to read element %d: %v", tc+1, i+1, err)
			}
		}
		testData[tc] = arr
	}
	return t, testData, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(200)...)

	for idx, tc := range tests {
		t, arrays, err := extractTestData(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid test input: %v\n", idx+1, err)
			os.Exit(1)
		}

		expOut, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expCounts, _, err := parseOutput(expOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotCounts, gotAppends, err := parseOutput(gotOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < t; caseIdx++ {
			if gotCounts[caseIdx] != expCounts[caseIdx] {
				fmt.Fprintf(os.Stderr, "case %d test %d: append count mismatch: expected %d got %d\n", idx+1, caseIdx+1, expCounts[caseIdx], gotCounts[caseIdx])
				os.Exit(1)
			}
			if gotCounts[caseIdx] == 0 {
				continue
			}
			if !validateSolution(arrays[caseIdx], gotAppends[caseIdx]) {
				fmt.Fprintf(os.Stderr, "case %d test %d: appended values do not form a good array\n", idx+1, caseIdx+1)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
