package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceE = "871E.go"
	refBinaryE = "ref871E.bin"
	totalTests = 80
)

type testCase struct {
	n int
	k int
	d [][]int
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

		if strings.TrimSpace(refOut) == "-1" {
			if strings.TrimSpace(candOut) != "-1" {
				fmt.Printf("test %d failed: expected -1, got\n%s\n", idx+1, candOut)
				printInput(input)
				os.Exit(1)
			}
			continue
		}

		if strings.TrimSpace(candOut) == "-1" {
			fmt.Printf("test %d failed: candidate reported -1 but solution exists\n", idx+1)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			os.Exit(1)
		}

		edges, err := parseEdges(candOut, tc.n)
		if err != nil {
			fmt.Printf("test %d failed: invalid candidate edges: %v\n", idx+1, err)
			printInput(input)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}

		if err := validateTree(tc, edges); err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			printInput(input)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
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

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i := 0; i < tc.k; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.d[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

type edge struct {
	u int
	v int
}

func parseEdges(out string, n int) ([]edge, error) {
	lines := strings.Fields(out)
	if len(lines)%2 != 0 {
		return nil, fmt.Errorf("expected pairs of vertices")
	}
	if len(lines)/2 != n-1 {
		return nil, fmt.Errorf("expected %d edges, got %d values", n-1, len(lines)/2)
	}
	edges := make([]edge, 0, n-1)
	for i := 0; i < len(lines); i += 2 {
		u, err := strconv.Atoi(lines[i])
		if err != nil {
			return nil, fmt.Errorf("invalid vertex %q", lines[i])
		}
		v, err := strconv.Atoi(lines[i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid vertex %q", lines[i+1])
		}
		if u < 1 || u > n || v < 1 || v > n || u == v {
			return nil, fmt.Errorf("invalid edge %d %d", u, v)
		}
		edges = append(edges, edge{u, v})
	}
	return edges, nil
}

func validateTree(tc testCase, edges []edge) error {
	if len(edges) != tc.n-1 {
		return fmt.Errorf("expected %d edges, got %d", tc.n-1, len(edges))
	}
	ds := newDSU(tc.n + 1)
	adj := make([][]int, tc.n+1)
	for _, e := range edges {
		if ds.find(e.u) == ds.find(e.v) {
			return fmt.Errorf("cycle detected")
		}
		ds.union(e.u, e.v)
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	root := 1
	visited := 0
	queue := []int{root}
	seen := make([]bool, tc.n+1)
	seen[root] = true
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		visited++
		for _, v := range adj[u] {
			if !seen[v] {
				seen[v] = true
				queue = append(queue, v)
			}
		}
	}
	if visited != tc.n {
		return fmt.Errorf("tree is not connected")
	}
	for i := 0; i < tc.k; i++ {
		dist := bfsDistances(adj, tc.n, i+1)
		if len(dist) != tc.n {
			return fmt.Errorf("distance computation failed")
		}
		for j := 0; j < tc.n; j++ {
			if dist[j] != tc.d[i][j] {
				return fmt.Errorf("distance mismatch for remembered vertex %d to vertex %d: expected %d, got %d",
					i+1, j+1, tc.d[i][j], dist[j])
			}
		}
	}
	return nil
}

func bfsDistances(adj [][]int, n int, start int) []int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt32
	}
	queue := []int{start}
	dist[start-1] = 0
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if dist[v-1] == math.MaxInt32 {
				dist[v-1] = dist[u-1] + 1
				queue = append(queue, v)
			}
		}
	}
	return dist
}

type dsu struct {
	parent []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &dsu{parent: p}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra != rb {
		d.parent[rb] = ra
	}
}

func generateTests() []testCase {
	tests := []testCase{
		buildLineTree(2, 1),
		buildLineTree(5, 2),
		buildStarTree(6, 3),
		buildRandomTree(8, 4, rand.New(rand.NewSource(1))),
		buildRandomTree(12, 5, rand.New(rand.NewSource(2))),
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		n := rnd.Intn(40) + 2
		k := rnd.Intn(min(5, n)) + 1
		tests = append(tests, buildRandomTree(n, k, rnd))
	}
	tests = append(tests,
		buildRandomTree(60, 5, rand.New(rand.NewSource(3))),
		buildRandomTree(80, 5, rand.New(rand.NewSource(4))),
		buildRandomTree(100, 6, rand.New(rand.NewSource(5))),
		buildRandomTree(200, 10, rand.New(rand.NewSource(6))),
		buildRandomTree(300, 10, rand.New(rand.NewSource(7))),
	)
	return tests
}

func buildLineTree(n, k int) testCase {
	edges := make([]edge, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, edge{i, i + 1})
	}
	return buildTestCase(n, k, edges)
}

func buildStarTree(n, k int) testCase {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, edge{1, i})
	}
	return buildTestCase(n, k, edges)
}

func buildRandomTree(n, k int, rnd *rand.Rand) testCase {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := rnd.Intn(i-1) + 1
		edges = append(edges, edge{parent, i})
	}
	return buildTestCase(n, k, edges)
}

func buildTestCase(n, k int, edges []edge) testCase {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	k = min(k, n)
	distances := make([][]int, k)
	for i := 0; i < k; i++ {
		start := i + 1
		distances[i] = bfsDistances(adj, n, start)
	}
	return testCase{n: n, k: k, d: distances}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printInput(in string) {
	fmt.Println("Input used:")
	fmt.Println(in)
}
