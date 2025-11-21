package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceE = "535E.go"
	refBinaryE = "ref535E.bin"
	totalTests = 80
)

type competitor struct {
	s int
	r int
}

type testCase struct {
	data []competitor
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}

		refAns, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Printf("test %d failed: expected %d indices, got %d\n", idx+1, len(refAns), len(candAns))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}

		refList := append([]int(nil), refAns...)
		candList := append([]int(nil), candAns...)
		sort.Ints(refList)
		sort.Ints(candList)

		for i := range refList {
			if refList[i] != candList[i] {
				fmt.Printf("test %d failed: mismatch at position %d (expected %d, got %d)\n", idx+1, i+1, refList[i], candList[i])
				printInput(input)
				fmt.Println("Reference output:")
				fmt.Println(refOut)
				fmt.Println("Candidate output:")
				fmt.Println(candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryE, refSourceE)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryE), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tc.data)))
	sb.WriteByte('\n')
	for _, comp := range tc.data {
		sb.WriteString(fmt.Sprintf("%d %d\n", comp.s, comp.r))
	}
	return []byte(sb.String())
}

func parseOutput(out string) ([]int, error) {
	fields := strings.Fields(out)
	res := make([]int, 0, len(fields))
	for _, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res = append(res, val)
	}
	return res, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{data: []competitor{{s: 1, r: 1}}},
		{data: []competitor{{1, 2}, {2, 1}}},
		{data: []competitor{{3, 3}, {3, 3}, {2, 4}, {4, 2}}},
		{data: []competitor{{5, 1}, {4, 2}, {3, 3}, {2, 4}, {1, 5}}},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(tests) < totalTests-4 {
		n := rnd.Intn(50) + 1
		data := make([]competitor, n)
		for i := 0; i < n; i++ {
			data[i] = competitor{
				s: rnd.Intn(10_000) + 1,
				r: rnd.Intn(10_000) + 1,
			}
		}
		tests = append(tests, testCase{data: data})
	}

	tests = append(tests, testCase{data: buildGradient(200)})
	tests = append(tests, testCase{data: buildGradient(2000)})
	tests = append(tests, testCase{data: buildClusters(5000)})
	tests = append(tests, testCase{data: buildLarge(200000)})

	return tests
}

func buildGradient(n int) []competitor {
	data := make([]competitor, n)
	for i := 0; i < n; i++ {
		data[i] = competitor{
			s: 10_000 - (i % 10_000),
			r: 1 + (i % 10_000),
		}
	}
	return data
}

func buildClusters(n int) []competitor {
	rnd := rand.New(rand.NewSource(2024))
	data := make([]competitor, n)
	for i := 0; i < n; i++ {
		baseS := rnd.Intn(5) * 2000
		baseR := rnd.Intn(5) * 2000
		data[i] = competitor{
			s: 1 + baseS + rnd.Intn(2000),
			r: 1 + baseR + rnd.Intn(2000),
		}
		if data[i].s > 10_000 {
			data[i].s = 10_000
		}
		if data[i].r > 10_000 {
			data[i].r = 10_000
		}
	}
	return data
}

func buildLarge(n int) []competitor {
	rnd := rand.New(rand.NewSource(424242))
	data := make([]competitor, n)
	for i := 0; i < n; i++ {
		data[i] = competitor{
			s: rnd.Intn(10_000) + 1,
			r: rnd.Intn(10_000) + 1,
		}
	}
	return data
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
