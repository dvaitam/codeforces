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
	wl [][2]int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1662H-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleH")
	cmd := exec.Command("go", "build", "-o", path, "1662H.go")
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
	return strings.TrimSpace(stdout.String()), nil
}

func parseLine(line string, expect int64) ([]int64, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty line")
	}
	k, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid k: %v", err)
	}
	if int64(len(fields)-1) != k {
		return nil, fmt.Errorf("reported k=%d but got %d values", k, len(fields)-1)
	}
	res := make([]int64, k)
	var prev int64 = -1
	for i := int64(0); i < k; i++ {
		val, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[i+1])
		}
		if val < 1 {
			return nil, fmt.Errorf("value %d not positive", val)
		}
		if val <= prev {
			return nil, fmt.Errorf("values not strictly increasing at position %d", i+1)
		}
		res[i] = val
		prev = val
	}
	return res, nil
}

func parseOutput(out string, count int) ([][]int64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != count {
		return nil, fmt.Errorf("expected %d lines, got %d", count, len(lines))
	}
	result := make([][]int64, count)
	for i := 0; i < count; i++ {
		vals, err := parseLine(lines[i], 0)
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", i+1, err)
		}
		result[i] = vals
	}
	return result, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.wl)))
	for _, pair := range tc.wl {
		sb.WriteString(fmt.Sprintf("%d %d\n", pair[0], pair[1]))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{wl: [][2]int64{{3, 5}}},
		{wl: [][2]int64{{3, 5}, {4, 4}, {5, 3}}},
		{wl: [][2]int64{{10, 10}, {7, 4}}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(5) + 1
	wl := make([][2]int64, t)
	for i := 0; i < t; i++ {
		w := int64(rng.Intn(1000) + 3)
		l := int64(rng.Intn(1000) + 3)
		wl[i] = [2]int64{w, l}
	}
	return testCase{wl: wl}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expVals, err := parseOutput(expOut, len(tc.wl))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, len(tc.wl))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		for i := range expVals {
			if len(expVals[i]) != len(gotVals[i]) {
				fmt.Fprintf(os.Stderr, "test %d case %d: expected count %d got %d\ninput:\n%s\n", idx+1, i+1, len(expVals[i]), len(gotVals[i]), input)
				os.Exit(1)
			}
			for j := range expVals[i] {
				if expVals[i][j] != gotVals[i][j] {
					fmt.Fprintf(os.Stderr, "test %d case %d: expected value %d got %d\ninput:\n%s\n", idx+1, i+1, expVals[i][j], gotVals[i][j], input)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
