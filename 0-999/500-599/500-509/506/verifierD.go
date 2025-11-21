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
	refSourceD = "506D.go"
	refBinaryD = "ref506D.bin"
	totalTests = 90
)

type query struct {
	u int
	v int
}

type edge struct {
	u int
	v int
	c int
}

type testCase struct {
	n     int
	m     int
	edges []edge
	qs    []query
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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

		refVals, err := parseOutput(refOut, len(tc.qs))
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, len(tc.qs))
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Printf("test %d failed: expected %d numbers, got %d\n", idx+1, len(refVals), len(candVals))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Printf("test %d failed at line %d: expected %d, got %d\n", idx+1, i+1, refVals[i], candVals[i])
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
	cmd := exec.Command("go", "build", "-o", refBinaryD, refSourceD)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryD), nil
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
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.c))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.qs)))
	for _, q := range tc.qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.u, q.v))
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
		buildManualTest(3, []edge{{1, 2, 1}, {2, 3, 2}}, []query{{1, 2}, {1, 3}, {2, 3}}),
		buildManualTest(4, []edge{{1, 2, 1}, {2, 3, 1}, {3, 4, 2}, {1, 4, 3}}, []query{{1, 3}, {1, 4}, {2, 4}}),
		buildManualTest(5, []edge{{1, 2, 1}, {2, 3, 2}, {3, 4, 1}, {4, 5, 2}, {1, 5, 3}}, []query{{1, 2}, {2, 3}, {1, 5}, {3, 5}}),
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(tests) < totalTests-5 {
		n := rnd.Intn(40) + 2
		m := rnd.Intn(n*(n-1)/2) + n
		maxColor := m
		edges := make([]edge, 0, m)
		used := make(map[[3]int]bool)
		for len(edges) < m {
			u := rnd.Intn(n) + 1
			v := rnd.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			c := rnd.Intn(maxColor) + 1
			key := [3]int{u, v, c}
			if used[key] {
				continue
			}
			used[key] = true
			edges = append(edges, edge{u: u, v: v, c: c})
		}
		qCnt := rnd.Intn(500) + 1
		queries := make([]query, qCnt)
		for i := 0; i < qCnt; i++ {
			u := rnd.Intn(n) + 1
			v := rnd.Intn(n) + 1
			for u == v {
				v = rnd.Intn(n) + 1
			}
			queries[i] = query{u: u, v: v}
		}
		tests = append(tests, testCase{n: n, m: m, edges: edges, qs: queries})
	}

	tests = append(tests, buildStarGraph(200, 400, 800))
	tests = append(tests, buildChainGraph(300, 300, 900))
	tests = append(tests, buildRandomDense(400, 1000, 1500))
	tests = append(tests, buildColorFocused(150, 600, 2000))
	tests = append(tests, buildHeavyColorCase(200, 2000, 2000))

	return tests
}

func buildManualTest(n int, edges []edge, qs []query) testCase {
	return testCase{
		n:     n,
		m:     len(edges),
		edges: append([]edge(nil), edges...),
		qs:    append([]query(nil), qs...),
	}
}

func buildStarGraph(n, colors, queries int) testCase {
	edges := make([]edge, 0, colors)
	for c := 1; c <= colors; c++ {
		u := 1
		v := (c % (n - 1)) + 2
		edges = append(edges, edge{u: u, v: v, c: c})
	}
	qs := make([]query, queries)
	for i := 0; i < queries; i++ {
		u := 1
		v := (i % (n - 1)) + 2
		qs[i] = query{u: u, v: v}
	}
	return testCase{n: n, m: len(edges), edges: edges, qs: qs}
}

func buildChainGraph(n, colors, queries int) testCase {
	edges := make([]edge, 0, n-1)
	color := 1
	for i := 1; i < n; i++ {
		edges = append(edges, edge{u: i, v: i + 1, c: color})
		color++
		if color > colors {
			color = 1
		}
	}
	qs := make([]query, queries)
	for i := 0; i < queries; i++ {
		u := 1
		v := n
		qs[i] = query{u: u, v: v}
	}
	return testCase{n: n, m: len(edges), edges: edges, qs: qs}
}

func buildRandomDense(n, m, q int) testCase {
	rnd := rand.New(rand.NewSource(12345))
	edges := make([]edge, 0, m)
	used := make(map[[3]int]bool)
	for len(edges) < m {
		u := rnd.Intn(n) + 1
		v := rnd.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		c := rnd.Intn(m) + 1
		key := [3]int{u, v, c}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, edge{u: u, v: v, c: c})
	}
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		u := rnd.Intn(n) + 1
		v := rnd.Intn(n) + 1
		for u == v {
			v = rnd.Intn(n) + 1
		}
		qs[i] = query{u: u, v: v}
	}
	return testCase{n: n, m: len(edges), edges: edges, qs: qs}
}

func buildColorFocused(n, colors, q int) testCase {
	rnd := rand.New(rand.NewSource(2024))
	edges := make([]edge, 0, colors)
	for c := 1; c <= colors; c++ {
		base := rnd.Intn(3) + 1
		u := rnd.Intn(n-1) + 1
		v := rnd.Intn(n-u) + u + 1
		if u > v {
			u, v = v, u
		}
		if base == 1 {
			v = u + 1
			if v > n {
				v = n
			}
		}
		edges = append(edges, edge{u: u, v: v, c: (c % colors) + 1})
	}
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		u := rnd.Intn(n) + 1
		v := rnd.Intn(n) + 1
		for u == v {
			v = rnd.Intn(n) + 1
		}
		qs[i] = query{u: u, v: v}
	}
	return testCase{n: n, m: len(edges), edges: edges, qs: qs}
}

func buildHeavyColorCase(n, colors, q int) testCase {
	edges := make([]edge, 0, colors)
	for c := 1; c <= colors; c++ {
		u := ((c - 1) % (n - 1)) + 1
		v := u + 1
		if v > n {
			v = n
		}
		if u == v {
			v = (v % n) + 1
		}
		if u > v {
			u, v = v, u
		}
		edges = append(edges, edge{u: u, v: v, c: c})
	}
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		u := (i % n) + 1
		v := ((i + n/2) % n) + 1
		if u == v {
			v = (v % n) + 1
		}
		qs[i] = query{u: u, v: v}
	}
	return testCase{n: n, m: len(edges), edges: edges, qs: qs}
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
