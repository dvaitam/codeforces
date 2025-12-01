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

const refSource = "./2063C.go"

type treeCase struct {
	n     int
	edges [][2]int
}

type testSuite struct {
	name  string
	cases []treeCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	suites := buildTests()

	for idx, suite := range suites {
		input := buildInput(suite.cases)
		expected := evaluate(suite.cases)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut, len(suite.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, refOut)
			os.Exit(1)
		}
		if !equalAnswers(refAns, expected) {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %d (%s)\ninput:\n%sreference:\n%s", idx+1, suite.name, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, len(suite.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, candOut)
			os.Exit(1)
		}
		if !equalAnswers(candAns, expected) {
			fmt.Fprintf(os.Stderr, "candidate mismatch on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, suite.name, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d test suites passed\n", len(suites))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2063C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2063C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, expected int) ([]int, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int, expected)
	for i, tok := range fields {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func equalAnswers(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildInput(cases []treeCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", c.n))
		for _, e := range c.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func evaluate(cases []treeCase) []int {
	ans := make([]int, len(cases))
	for i, c := range cases {
		ans[i] = solveExact(c)
	}
	return ans
}

func solveExact(c treeCase) int {
	n := c.n
	deg := make([]int, n+1)
	adj := make(map[uint64]struct{}, len(c.edges))
	for _, e := range c.edges {
		u, v := e[0], e[1]
		deg[u]++
		deg[v]++
		adj[edgeKey(u, v)] = struct{}{}
	}
	best := 0
	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			edge := 0
			if _, ok := adj[edgeKey(u, v)]; ok {
				edge = 1
			}
			val := deg[u] + deg[v] - 1 - edge
			if val > best {
				best = val
			}
		}
	}
	return best
}

func edgeKey(u, v int) uint64 {
	if u > v {
		u, v = v, u
	}
	return (uint64(u) << 32) | uint64(v)
}

func buildTests() []testSuite {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	deterministic := []treeCase{
		// Small path
		{
			n: 2,
			edges: [][2]int{
				{1, 2},
			},
		},
		{
			n: 4,
			edges: [][2]int{
				{1, 2}, {2, 3}, {3, 4},
			},
		},
		// Star to check use of high degree center
		{
			n: 5,
			edges: [][2]int{
				{1, 2}, {1, 3}, {1, 4}, {1, 5},
			},
		},
		// Two adjacent high-degree vertices
		{
			n: 6,
			edges: [][2]int{
				{1, 2}, {1, 3}, {1, 4}, {2, 5}, {2, 6},
			},
		},
		// Balanced branching
		{
			n: 7,
			edges: [][2]int{
				{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}, {3, 7},
			},
		},
	}

	var randomCases []treeCase
	for i := 0; i < 6; i++ {
		n := 10 + i*5
		randomCases = append(randomCases, randomTree(rng, n))
	}

	combinedSmall := testSuite{name: "deterministic-small", cases: deterministic}
	combinedRandom := testSuite{name: "randomized", cases: randomCases}
	mixed := testSuite{name: "mixed", cases: append([]treeCase{deterministic[0], deterministic[3]}, randomCases[:2]...)}

	return []testSuite{combinedSmall, combinedRandom, mixed}
}

func randomTree(rng *rand.Rand, n int) treeCase {
	if n < 2 {
		n = 2
	}
	prufer := make([]int, n-2)
	for i := range prufer {
		prufer[i] = rng.Intn(n) + 1
	}
	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = 1
	}
	for _, x := range prufer {
		deg[x]++
	}
	edges := make([][2]int, 0, n-1)
	leaf := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			leaf = append(leaf, i)
		}
	}
	sortInts(leaf)
	for _, x := range prufer {
		u := leaf[0]
		leaf = leaf[1:]
		edges = append(edges, [2]int{u, x})
		deg[u]--
		deg[x]--
		if deg[x] == 1 {
			insertSorted(&leaf, x)
		}
	}
	// last edge
	var last []int
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			last = append(last, i)
		}
	}
	edges = append(edges, [2]int{last[0], last[1]})
	return treeCase{n: n, edges: edges}
}

func insertSorted(arr *[]int, val int) {
	a := *arr
	i := len(a)
	for i > 0 && a[i-1] > val {
		i--
	}
	a = append(a, 0)
	copy(a[i+1:], a[i:])
	a[i] = val
	*arr = a
}

func sortInts(a []int) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		val := a[i]
		j := i - 1
		for j >= 0 && a[j] > val {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = val
	}
}
