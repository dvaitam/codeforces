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
	refSource = "./1916F.go"
	maxEdges  = 5000
)

type testCase struct {
	n1, n2 int
	edges  [][2]int
	adj    [][]int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	tests := buildTests(rng)

	for idx, tc := range tests {
		input := tc.toInput()
		if err := checkProgram(refBin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if err := checkProgram(candidate, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	outPath := "./ref_1916F.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func checkProgram(bin string, input string, tc testCase) error {
	out, err := runProgram(bin, input)
	if err != nil {
		return err
	}
	return validateOutput(out, tc)
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests(rng *rand.Rand) []testCase {
	var tests []testCase

	tests = append(tests,
		makeCycleTest(1, 1, 0, rng),
		makeCycleTest(1, 2, 1, rng),
		makeCycleTest(2, 3, 1, rng),
		makeCycleTest(3, 4, 2, rng),
	)

	for i := 0; i < 40; i++ {
		n1 := rng.Intn(4) + 1
		n2 := rng.Intn(4) + 1
		extra := rng.Intn(3)
		tests = append(tests, makeCycleTest(n1, n2, extra, rng))
	}

	for i := 0; i < 25; i++ {
		total := rng.Intn(40) + 10
		n1 := rng.Intn(total-1) + 1
		n2 := total - n1
		extra := rng.Intn(total)
		tests = append(tests, makeCycleTest(n1, n2, extra, rng))
	}

	for _, total := range []int{50, 100, 200, 500, 1000, 1500, 2000} {
		n1 := rng.Intn(total-1) + 1
		n2 := total - n1
		tests = append(tests, makeCycleTest(n1, n2, total, rng))
	}

	return tests
}

func makeCycleTest(n1, n2, extra int, rng *rand.Rand) testCase {
	n := n1 + n2
	edges := buildCycleEdges(n)
	seen := make(map[int]struct{}, len(edges))
	for _, e := range edges {
		key := edgeKey(e[0], e[1])
		seen[key] = struct{}{}
	}
	target := len(edges) + extra
	if target > maxEdges {
		target = maxEdges
	}
	attempts := 0
	for len(edges) < target && attempts < extra*5+100 {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		if addEdge(&edges, u, v, seen) {
			continue
		}
		attempts++
	}
	adj := buildAdj(n, edges)
	return testCase{n1: n1, n2: n2, edges: edges, adj: adj}
}

func buildCycleEdges(n int) [][2]int {
	var edges [][2]int
	seen := make(map[int]struct{})
	for i := 1; i <= n; i++ {
		u := i
		v := i%n + 1
		addEdge(&edges, u, v, seen)
	}
	return edges
}

func addEdge(edges *[][2]int, u, v int, seen map[int]struct{}) bool {
	if len(*edges) >= maxEdges || u == v {
		return false
	}
	if u > v {
		u, v = v, u
	}
	key := edgeKey(u, v)
	if _, ok := seen[key]; ok {
		return false
	}
	seen[key] = struct{}{}
	*edges = append(*edges, [2]int{u, v})
	return true
}

func edgeKey(u, v int) int {
	if u > v {
		u, v = v, u
	}
	return u*10000 + v
}

func buildAdj(n int, edges [][2]int) [][]int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	return adj
}

func (tc testCase) toInput() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n1, tc.n2, len(tc.edges)))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func validateOutput(out string, tc testCase) error {
	tokens := strings.Fields(out)
	expected := tc.n1 + tc.n2
	if len(tokens) != expected {
		return fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	n := tc.n1 + tc.n2
	group1 := make([]int, tc.n1)
	group2 := make([]int, tc.n2)
	used := make([]bool, n+1)

	for i := 0; i < tc.n1; i++ {
		val, err := strconv.Atoi(tokens[i])
		if err != nil {
			return fmt.Errorf("invalid integer in first group: %v", err)
		}
		if val < 1 || val > n {
			return fmt.Errorf("value %d out of range in first group", val)
		}
		if used[val] {
			return fmt.Errorf("value %d repeated", val)
		}
		used[val] = true
		group1[i] = val
	}
	for i := 0; i < tc.n2; i++ {
		val, err := strconv.Atoi(tokens[tc.n1+i])
		if err != nil {
			return fmt.Errorf("invalid integer in second group: %v", err)
		}
		if val < 1 || val > n {
			return fmt.Errorf("value %d out of range in second group", val)
		}
		if used[val] {
			return fmt.Errorf("value %d repeated", val)
		}
		used[val] = true
		group2[i] = val
	}
	for i := 1; i <= n; i++ {
		if !used[i] {
			return fmt.Errorf("value %d missing from both groups", i)
		}
	}
	if !isConnected(group1, tc.adj) {
		return fmt.Errorf("first group is not connected")
	}
	if !isConnected(group2, tc.adj) {
		return fmt.Errorf("second group is not connected")
	}
	return nil
}

func isConnected(nodes []int, adj [][]int) bool {
	if len(nodes) == 0 {
		return false
	}
	in := make([]bool, len(adj))
	for _, v := range nodes {
		in[v] = true
	}
	visited := make([]bool, len(adj))
	queue := []int{nodes[0]}
	visited[nodes[0]] = true
	count := 0
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		count++
		for _, w := range adj[u] {
			if in[w] && !visited[w] {
				visited[w] = true
				queue = append(queue, w)
			}
		}
	}
	return count == len(nodes)
}
