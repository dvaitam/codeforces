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
	refSourceB1 = "324B1.go"
	refBinaryB1 = "ref324B1.bin"
	totalTests  = 70
)

type query struct {
	typ int
	x   int
	y   int
}

type testCase struct {
	n       int
	perm    []int
	queries []query
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
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

		refVals, err := parseOutput(refOut, countTypeOne(tc.queries))
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, countTypeOne(tc.queries))
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Printf("test %d failed: expected %d answers, got %d\n", idx+1, len(refVals), len(candVals))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Printf("test %d failed at answer %d: expected %d, got %d\n", idx+1, i+1, refVals[i], candVals[i])
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
	cmd := exec.Command("go", "build", "-o", refBinaryB1, refSourceB1)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryB1), nil
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
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(len(tc.queries)))
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", q.typ, q.x, q.y))
	}
	return []byte(sb.String())
}

func parseOutput(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int, expected)
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
		buildTestCase([]int{1, 2}, []query{{typ: 1, x: 1, y: 2}}),
		buildTestCase([]int{2, 1, 3}, []query{{typ: 1, x: 1, y: 3}, {typ: 2, x: 1, y: 2}, {typ: 1, x: 1, y: 3}}),
		buildTestCase([]int{3, 1, 2, 4}, []query{{typ: 1, x: 1, y: 4}, {typ: 1, x: 2, y: 4}}),
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-3 {
		n := rnd.Intn(90) + 10
		perm := randPermutation(rnd, n)
		qCount := rnd.Intn(200) + 50
		queries := randomQueries(rnd, n, qCount)
		tests = append(tests, buildTestCase(perm, queries))
	}

	tests = append(tests, buildTestCase(randPermutation(rand.New(rand.NewSource(1)), 100), randomQueries(rand.New(rand.NewSource(2)), 100, 1000)))
	tests = append(tests, buildTestCase(randPermutation(rand.New(rand.NewSource(3)), 50), randomQueries(rand.New(rand.NewSource(4)), 50, 5000)))
	tests = append(tests, buildTestCase([]int{4, 3, 2, 1}, []query{
		{typ: 1, x: 1, y: 4},
		{typ: 2, x: 1, y: 4},
		{typ: 1, x: 1, y: 4},
		{typ: 2, x: 2, y: 3},
		{typ: 1, x: 1, y: 4},
	}))

	return tests
}

func buildTestCase(perm []int, queries []query) testCase {
	return testCase{
		n:       len(perm),
		perm:    append([]int(nil), perm...),
		queries: append([]query(nil), queries...),
	}
}

func randPermutation(rnd *rand.Rand, n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rnd.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}

func randomQueries(rnd *rand.Rand, n, count int) []query {
	queries := make([]query, count)
	for i := 0; i < count; i++ {
		typ := rnd.Intn(2) + 1
		x := rnd.Intn(n-1) + 1
		y := rnd.Intn(n-x) + x + 1
		queries[i] = query{typ: typ, x: x, y: y}
	}
	return queries
}

func countTypeOne(queries []query) int {
	total := 0
	for _, q := range queries {
		if q.typ == 1 {
			total++
		}
	}
	return total
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
