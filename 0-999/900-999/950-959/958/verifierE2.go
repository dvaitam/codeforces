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

const (
	refSourceE2 = "958E2.go"
	refBinaryE2 = "ref958E2.bin"
	totalTests  = 70
)

type testCase struct {
	k     int
	n     int
	times []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Printf("test %d failed: expected %d, got %d\n", idx+1, refVal, candVal)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryE2, refSourceE2)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryE2), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.k, tc.n))
	for i, t := range tc.times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(t, 10))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func generateTests() []testCase {
	tests := []testCase{
		{2, 4, []int64{1, 4, 6, 7}},
		{3, 6, []int64{1, 2, 3, 4, 5, 6}},
		{4, 8, []int64{1, 3, 4, 5, 14, 15, 23, 25}},
		{2, 5, []int64{10, 1, 7, 4, 20}},
		{5, 12, []int64{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60}},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		n := rnd.Intn(200) + 4
		k := rnd.Intn(min(10, n/2)) + 2
		times := make([]int64, n)
		used := make(map[int64]bool)
		for i := 0; i < n; i++ {
			for {
				val := rnd.Int63n(1_000_000_000) + 1
				if !used[val] {
					used[val] = true
					times[i] = val
					break
				}
			}
		}
		tests = append(tests, testCase{k: k, n: n, times: times})
	}

	tests = append(tests,
		randomLarge(500, 50, rand.New(rand.NewSource(1))),
		randomLarge(1000, 100, rand.New(rand.NewSource(2))),
		randomLarge(5000, 300, rand.New(rand.NewSource(3))),
		randomLarge(20000, 1000, rand.New(rand.NewSource(4))),
		randomLarge(50000, 2500, rand.New(rand.NewSource(5))),
	)

	return tests
}

func randomLarge(n, k int, rnd *rand.Rand) testCase {
	times := make([]int64, n)
	used := make(map[int64]bool, n)
	for i := 0; i < n; i++ {
		for {
			val := rnd.Int63n(1_000_000_000) + 1
			if !used[val] {
				used[val] = true
				times[i] = val
				break
			}
		}
	}
	return testCase{k: k, n: n, times: times}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
