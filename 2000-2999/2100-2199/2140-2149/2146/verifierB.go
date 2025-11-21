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
	refSource        = "2146B.go"
	tempOraclePrefix = "oracle-2146B-"
	randomTests      = 60
)

type testCase struct {
	n    int
	m    int
	sets [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestsCases(rng, randomTests)...)

	for idx, tc := range tests {
		input := formatInput(tc)
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
		exp := normalize(expOut)
		got := normalize(gotOut)
		if exp == "" || got == "" {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", idx+1)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", idx+1, exp, got)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(gotOut)
			os.Exit(1)
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
	outPath := filepath.Join(tmpDir, "oracleB")
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

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.m)
	for _, s := range tc.sets {
		fmt.Fprintf(&sb, "%d", len(s))
		for _, v := range s {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func normalize(out string) string {
	out = strings.TrimSpace(strings.ToLower(out))
	if out == "yes" {
		return "yes"
	}
	if out == "no" {
		return "no"
	}
	return ""
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, m: 2, sets: [][]int{{1}, {2}}},
		{n: 3, m: 3, sets: [][]int{{1, 2}, {2, 3}, {1, 3}}},
		{n: 4, m: 4, sets: [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 4}}},
		{n: 5, m: 5, sets: [][]int{{1}, {2}, {3}, {4}, {5}}},
		{n: 3, m: 4, sets: [][]int{{1, 2}, {2, 3}, {3, 4}}},
	}
}

func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		n := rng.Intn(6) + 2
		m := rng.Intn(6) + 2
		sets := make([][]int, n)
		used := make([]bool, m+1)
		for i := 0; i < n; i++ {
			size := rng.Intn(m) + 1
			elems := randSubset(rng, m, size)
			for _, v := range elems {
				used[v] = true
			}
			sets[i] = elems
		}
		valid := true
		for v := 1; v <= m; v++ {
			if !used[v] {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		tests = append(tests, testCase{n: n, m: m, sets: sets})
	}
	return tests
}

func randSubset(rng *rand.Rand, m, size int) []int {
	avail := rand.Perm(m)
	elems := make([]int, size)
	for i := 0; i < size; i++ {
		elems[i] = avail[i] + 1
	}
	sortInts(elems)
	return elems
}

func sortInts(a []int) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && a[j] > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}
