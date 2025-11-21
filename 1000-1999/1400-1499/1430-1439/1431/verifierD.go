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
	a [][]int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1431D-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", path, "1431D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string, input string) (string, error) {
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

func parseOutput(out string, counts []int) ([][]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(counts) {
		return nil, fmt.Errorf("expected %d lines, got %d", len(counts), len(lines))
	}
	result := make([][]int, len(counts))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != counts[i] {
			return nil, fmt.Errorf("test case %d: expected %d numbers, got %d", i+1, counts[i], len(fields))
		}
		order := make([]int, counts[i])
		seen := make([]bool, counts[i])
		for j, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("test case %d: invalid integer %q", i+1, f)
			}
			if val < 1 || val > counts[i] {
				return nil, fmt.Errorf("test case %d: index %d out of range", i+1, val)
			}
			if seen[val-1] {
				return nil, fmt.Errorf("test case %d: duplicate index %d", i+1, val)
			}
			seen[val-1] = true
			order[j] = val - 1
		}
		result[i] = order
	}
	return result, nil
}

func simulate(order []int, a []int) int {
	usedMarkers := 1
	current := 0
	for _, idx := range order {
		if current >= a[idx] {
			usedMarkers++
			current = 1
		} else {
			current++
		}
	}
	return usedMarkers
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.a)))
	for _, arr := range tc.a {
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{a: [][]int{{1}}},
		{a: [][]int{{1, 2}, {2, 1}}},
		{a: [][]int{{3, 1, 2}, {1, 1, 1}}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(5) + 1
	cases := make([][]int, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(n) + 1
		}
		cases[i] = arr
	}
	return testCase{a: cases}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		counts := make([]int, len(tc.a))
		for i, arr := range tc.a {
			counts[i] = len(arr)
		}

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expOrders, err := parseOutput(expOut, counts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}
		expBest := make([]int, len(expOrders))
		for i, arr := range expOrders {
			expBest[i] = simulate(arr, tc.a[i])
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotOrders, err := parseOutput(gotOut, counts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}
		for i, arr := range gotOrders {
			score := simulate(arr, tc.a[i])
			if score != expBest[i] {
				fmt.Fprintf(os.Stderr, "test %d case %d: expected score %d got %d\ninput:\n%s\n", idx+1, i+1, expBest[i], score, input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
