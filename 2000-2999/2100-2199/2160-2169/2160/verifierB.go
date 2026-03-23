package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	randomTests = 80
	maxN        = 200
)

type testCase struct {
	n int
	b []int64
}

// Embedded solver (matches 2160B.go logic).
func solve(tc testCase) string {
	n := tc.n
	b := tc.b

	a := make([]int, n+1)
	nextVal := 1
	var prev int64

	for i := 1; i <= n; i++ {
		g := b[i] - prev
		prev = b[i]
		l := i - int(g)
		if l <= 0 {
			a[i] = nextVal
			nextVal++
		} else {
			a[i] = a[l]
		}
	}

	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestsCases(rng, randomTests)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		expOut := solve(tc)

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expVals := parseInts(expOut)
		gotVals := parseInts(gotOut)
		if expVals == nil || gotVals == nil || len(expVals) != tc.n || len(gotVals) != tc.n {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", idx+1)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "test %d mismatch at position %d: expected %d got %d\n", idx+1, i+1, expVals[i], gotVals[i])
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

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", tc.n)
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(tc.b[i], 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, b: []int64{0, 1}},
		{n: 3, b: []int64{0, 1, 3, 6}},
		{n: 3, b: []int64{0, 1, 3, 5}},
		{n: 3, b: []int64{0, 1, 3, 4}},
		{n: 4, b: []int64{0, 1, 2, 3, 7}},
	}
}

// Generate valid test cases by constructing array a first, then computing b.
func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		n := rng.Intn(maxN-1) + 1
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			a[i] = rng.Intn(n) + 1
		}
		// compute b from a
		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			var sum int64
			seen := make(map[int]bool)
			for j := i; j >= 1; j-- {
				seen[a[j]] = true
				sum += int64(len(seen))
			}
			b[i] = sum
		}
		tests = append(tests, testCase{n: n, b: b})
	}
	return tests
}
