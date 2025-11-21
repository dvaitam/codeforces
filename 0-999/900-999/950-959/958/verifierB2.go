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
	refSourceB2 = "958B2.go"
	refBinaryB2 = "ref958B2.bin"
	totalTests  = 80
)

type testCase struct {
	n     int
	edges [][2]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
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
		refVals, err := parseOutput(refOut, tc.n)
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
		candVals, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Printf("test %d failed: expected %d values, got %d\n", idx+1, len(refVals), len(candVals))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Printf("test %d failed at position %d: expected %d, got %d\n", idx+1, i+1, refVals[i], candVals[i])
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
	cmd := exec.Command("go", "build", "-o", refBinaryB2, refSourceB2)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryB2), nil
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

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return []byte(sb.String())
}

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	res := make([]int, n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildStar(3),
		buildLine(5),
		buildLine(10),
		buildRandom(20, rand.New(rand.NewSource(1))),
		buildRandom(40, rand.New(rand.NewSource(2))),
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		n := rnd.Intn(80) + 2
		tests = append(tests, buildRandom(n, rnd))
	}
	tests = append(tests,
		buildRandom(200, rand.New(rand.NewSource(3))),
		buildRandom(500, rand.New(rand.NewSource(4))),
		buildRandom(1000, rand.New(rand.NewSource(5))),
		buildRandom(5000, rand.New(rand.NewSource(6))),
		buildRandom(10000, rand.New(rand.NewSource(7))),
	)
	return tests
}

func buildStar(n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return testCase{n: n, edges: edges}
}

func buildLine(n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return testCase{n: n, edges: edges}
}

func buildRandom(n int, rnd *rand.Rand) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rnd.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return testCase{n: n, edges: edges}
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
