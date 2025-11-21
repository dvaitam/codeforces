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
	refSource        = "2155C.go"
	tempOraclePrefix = "oracle-2155C-"
	randomTests      = 80
	maxN             = 200
	maxValue         = 1000
)

type testCase struct {
	n int
	a []int
}

type testGroup struct {
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicGroups()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomGroups(rng, randomTests)...)

	for idx, tg := range tests {
		input := formatGroupInput(tg)
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expVals := parseInts(expOut)
		gotVals := parseInts(gotOut)
		if expVals == nil || gotVals == nil || len(expVals) != len(tg.cases) || len(gotVals) != len(tg.cases) {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", idx+1)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
		for i := range expVals {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "test %d case %d failed: expected %d got %d\n", idx+1, i+1, expVals[i], gotVals[i])
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Expected:")
				fmt.Print(expOut)
				fmt.Println("Got:")
				fmt.Print(gotOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func formatGroupInput(tg testGroup) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tg.cases))
	for _, tc := range tg.cases {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseInts(out string) []int64 {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil
	}
	res := make([]int64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil
		}
		res[i] = v
	}
	return res
}

func deterministicGroups() []testGroup {
	return []testGroup{
		{
			cases: []testCase{
				{n: 1, a: []int{5}},
				{n: 2, a: []int{1, 2}},
			},
		},
		{
			cases: []testCase{
				{n: 3, a: []int{3, 1, 2}},
				{n: 4, a: []int{5, -1, 0, 2}},
				{n: 5, a: []int{-2, -2, -2, -2, -2}},
			},
		},
	}
}

func randomGroups(rng *rand.Rand, count int) []testGroup {
	groups := make([]testGroup, 0, count)
	for len(groups) < count {
		tcCount := rng.Intn(5) + 1
		group := testGroup{cases: make([]testCase, tcCount)}
		totalN := 0
		for i := 0; i < tcCount; i++ {
			n := rng.Intn(maxN-1) + 1
			totalN += n
			values := make([]int, n)
			for j := 0; j < n; j++ {
				values[j] = rng.Intn(maxValue*2+1) - maxValue
			}
			group.cases[i] = testCase{n: n, a: values}
		}
		if totalN > maxN*2 {
			continue
		}
		groups = append(groups, group)
	}
	return groups
}
