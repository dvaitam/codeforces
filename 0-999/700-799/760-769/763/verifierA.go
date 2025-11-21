package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	n      int
	edges  [][2]int
	colors []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(output string, tc testCase) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	answer := strings.TrimSpace(lines[0])
	if strings.ToUpper(answer) == "NO" {
		if hasValidRoot(tc.n, tc.edges, tc.colors) {
			return fmt.Errorf("claimed NO but a valid root exists")
		}
		return nil
	}
	if strings.ToUpper(answer) != "YES" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if len(lines) < 2 {
		return fmt.Errorf("missing second line with vertex index")
	}
	rootLine := strings.TrimSpace(lines[1])
	root, err := strconv.Atoi(rootLine)
	if err != nil {
		return fmt.Errorf("invalid vertex index %q", rootLine)
	}
	if root < 1 || root > tc.n {
		return fmt.Errorf("vertex index %d out of range", root)
	}
	if !isValidRoot(root, tc.edges, tc.colors) {
		return fmt.Errorf("vertex %d does not satisfy the condition", root)
	}
	return nil
}

func isValidRoot(root int, edges [][2]int, colors []int) bool {
	for _, e := range edges {
		u, v := e[0], e[1]
		if colors[u] != colors[v] && root != u && root != v {
			return false
		}
	}
	return true
}

func hasValidRoot(n int, edges [][2]int, colors []int) bool {
	diffEdge := -1
	for i, e := range edges {
		if colors[e[0]] != colors[e[1]] {
			diffEdge = i
			break
		}
	}
	if diffEdge == -1 {
		return true
	}
	c1 := edges[diffEdge][0]
	c2 := edges[diffEdge][1]
	return isValidRoot(c1, edges, colors) || isValidRoot(c2, edges, colors)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []testCase {
	rand.Seed(42)
	var tests []testCase
	tests = append(tests, makeTest([][2]int{{1, 2}}, []int{0, 1, 1}))
	tests = append(tests, makeTest([][2]int{{1, 2}}, []int{0, 1, 2}))
	for i := 0; i < 200; i++ {
		n := rand.Intn(8) + 2
		edges := randomTree(n)
		colors := make([]int, n+1)
		for v := 1; v <= n; v++ {
			colors[v] = rand.Intn(3) + 1
		}
		tests = append(tests, makeTest(edges, colors))
	}
	return tests
}

func makeTest(edges [][2]int, colors []int) testCase {
	n := len(colors) - 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(colors[i]))
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		n:      n,
		edges:  edges,
		colors: colors,
	}
}

func randomTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rand.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	return edges
}
