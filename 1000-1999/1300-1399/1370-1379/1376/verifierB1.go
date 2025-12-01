package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type graph struct {
	n     int
	m     int
	edges [][2]int
}

type testCase struct {
	name  string
	graph graph
	input string
}

func makeInput(g graph) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", g.n, g.m)
	for _, e := range g.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(output string, n int) (int, []int, error) {
	lines := strings.Fields(output)
	if len(lines) < n+1 {
		return 0, nil, fmt.Errorf("expected at least %d integers, got %d", n+1, len(lines))
	}
	var k int
	if _, err := fmt.Sscan(lines[0], &k); err != nil {
		return 0, nil, fmt.Errorf("failed to parse k: %v", err)
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Sscan(lines[i+1], &arr[i]); err != nil {
			return 0, nil, fmt.Errorf("failed to parse x_%d: %v", i+1, err)
		}
		if arr[i] != 0 && arr[i] != 1 {
			return 0, nil, fmt.Errorf("x_%d is not 0 or 1", i+1)
		}
	}
	return k, arr, nil
}

func verifyIndependentSet(g graph, arr []int) error {
	selected := make(map[int]struct{})
	for i, v := range arr {
		if v == 1 {
			selected[i+1] = struct{}{}
		}
	}
	for _, e := range g.edges {
		if _, ok := selected[e[0]]; ok {
			if _, ok2 := selected[e[1]]; ok2 {
				return fmt.Errorf("edge (%d,%d) has both endpoints selected", e[0], e[1])
			}
		}
	}
	return nil
}

func randomGraph(name string, n int, edgeProb float64) testCase {
	var edges [][2]int
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if rand.Float64() < edgeProb {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	return testCase{
		name: name,
		graph: graph{
			n:     n,
			m:     len(edges),
			edges: edges,
		},
		input: makeInput(graph{n: n, m: len(edges), edges: edges}),
	}
}

func handcraftedTests() []testCase {
	// Simple graphs: single edge, empty graph, triangle, etc.
	tests := []testCase{}

	// Empty graph with 3 nodes
	tests = append(tests, testCase{
		name: "empty",
		graph: graph{
			n: 3,
			m: 0,
		},
		input: "3 0\n",
	})

	// Single edge
	tests = append(tests, testCase{
		name: "single_edge",
		graph: graph{
			n: 2,
			m: 1,
			edges: [][2]int{
				{1, 2},
			},
		},
		input: "2 1\n1 2\n",
	})

	// Triangle
	tests = append(tests, testCase{
		name: "triangle",
		graph: graph{
			n: 3,
			m: 3,
			edges: [][2]int{
				{1, 2}, {2, 3}, {1, 3},
			},
		},
		input: "3 3\n1 2\n2 3\n1 3\n",
	})

	// Path of 4 nodes
	tests = append(tests, testCase{
		name: "path4",
		graph: graph{
			n: 4,
			m: 3,
			edges: [][2]int{
				{1, 2}, {2, 3}, {3, 4},
			},
		},
		input: "4 3\n1 2\n2 3\n3 4\n",
	})

	return tests
}

func randomTests() []testCase {
	rand.Seed(time.Now().UnixNano())
	var tests []testCase
	sizes := []int{5, 6, 7, 8, 10}
	for idx, n := range sizes {
		tests = append(tests, randomGraph(fmt.Sprintf("rand_%d", idx+1), n, 0.3))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		output, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		k, arr, err := parseOutput(output, tc.graph.n)
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, output)
			os.Exit(1)
		}
		if err := verifyIndependentSet(tc.graph, arr); err != nil {
			fmt.Printf("test %d (%s) set is not independent: %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, output)
			os.Exit(1)
		}
		actualCount := 0
		for _, v := range arr {
			if v == 1 {
				actualCount++
			}
		}
		if actualCount != k {
			fmt.Printf("test %d (%s) mismatch in count: claimed %d but actual %d\ninput:\n%soutput:\n%s\n", idx+1, tc.name, k, actualCount, tc.input, output)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
