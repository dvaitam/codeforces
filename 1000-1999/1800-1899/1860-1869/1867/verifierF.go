package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fail("usage: verifierF /path/to/candidate")
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	numCases := 50
	for tc := 1; tc <= numCases; tc++ {
		inputData := genCase(rng)

		refOut, err := runProgram(exec.Command(refBin), []byte(inputData))
		if err != nil {
			fail("case %d: reference execution failed: %v", tc, err)
		}

		candOut, err := runProgram(commandFor(candidate), []byte(inputData))
		if err != nil {
			fail("case %d: candidate execution failed: %v", tc, err)
		}

		n := parseN(inputData)
		inputAdj := parseTree(inputData, n)

		refAdj, err := parseOutputTree(refOut, n)
		if err != nil {
			fail("case %d: invalid reference output: %v", tc, err)
		}
		candAdj, err := parseOutputTree(candOut, n)
		if err != nil {
			fail("case %d: invalid candidate output: %v\ncandidate output:\n%s", tc, err, candOut)
		}

		inputTypes := canonicalSubtreeTypes(inputAdj, n)
		refCost := countMatches(refAdj, n, inputTypes)
		candCost := countMatches(candAdj, n, inputTypes)

		if candCost > refCost {
			fail("case %d: candidate cost %d > reference cost %d\ninput:\n%scandidate output:\n%s",
				tc, candCost, refCost, inputData, candOut)
		}
	}

	fmt.Println("OK")
}

// genCase generates a random tree with n in [2..12].
func genCase(rng *rand.Rand) string {
	n := rng.Intn(11) + 2 // 2..12
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	// Shuffle edge order and randomly swap endpoints for realism.
	rng.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		a, b := e[0], e[1]
		if rng.Intn(2) == 0 {
			a, b = b, a
		}
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	return sb.String()
}

func parseN(input string) int {
	var n int
	fmt.Sscan(input, &n)
	return n
}

// parseTree parses the input tree into an adjacency list (1-indexed).
func parseTree(input string, n int) [][]int {
	adj := make([][]int, n+1)
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Scan() // skip first line (n)
	for i := 0; i < n-1; i++ {
		sc.Scan()
		var a, b int
		fmt.Sscan(sc.Text(), &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	return adj
}

// parseOutputTree parses the candidate/reference output as n-1 edges,
// validates it forms a tree on vertices 1..n.
func parseOutputTree(output string, n int) ([][]int, error) {
	adj := make([][]int, n+1)
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(output)))
	edgeCount := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var a, b int
		if _, err := fmt.Sscan(line, &a, &b); err != nil {
			return nil, fmt.Errorf("cannot parse edge %q: %v", line, err)
		}
		if a < 1 || a > n || b < 1 || b > n {
			return nil, fmt.Errorf("vertex out of range in edge %d %d", a, b)
		}
		if a == b {
			return nil, fmt.Errorf("self-loop %d %d", a, b)
		}
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
		edgeCount++
	}
	if edgeCount != n-1 {
		return nil, fmt.Errorf("expected %d edges, got %d", n-1, edgeCount)
	}
	// Check connectivity via BFS from vertex 1.
	visited := make([]bool, n+1)
	queue := []int{1}
	visited[1] = true
	cnt := 1
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, u := range adj[v] {
			if !visited[u] {
				visited[u] = true
				cnt++
				queue = append(queue, u)
			}
		}
	}
	if cnt != n {
		return nil, fmt.Errorf("tree is not connected (reached %d of %d vertices)", cnt, n)
	}
	return adj, nil
}

// canonicalSubtreeTypes computes a canonical string for each vertex's subtree
// when rooted at vertex 1 and returns the set of distinct canonical strings.
// A leaf has canonical form "()", a node with children c1..ck (sorted) is
// "(c1c2...ck)". This representation is independent of vertex numbering.
func canonicalSubtreeTypes(adj [][]int, n int) map[string]bool {
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	visited := make([]bool, n+1)
	visited[1] = true
	queue := []int{1}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, u := range adj[v] {
			if !visited[u] {
				visited[u] = true
				parent[u] = v
				queue = append(queue, u)
			}
		}
	}

	canon := make([]string, n+1)
	types := make(map[string]bool)

	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		childCanons := []string{}
		for _, u := range adj[v] {
			if u != parent[v] {
				childCanons = append(childCanons, canon[u])
			}
		}
		sort.Strings(childCanons)
		var sb strings.Builder
		sb.WriteByte('(')
		for _, c := range childCanons {
			sb.WriteString(c)
		}
		sb.WriteByte(')')
		canon[v] = sb.String()
		types[canon[v]] = true
	}
	return types
}

// countMatches counts how many vertices in the given tree (rooted at 1)
// have a canonical subtree type that also appears in inputTypes.
func countMatches(adj [][]int, n int, inputTypes map[string]bool) int {
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	visited := make([]bool, n+1)
	visited[1] = true
	queue := []int{1}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, u := range adj[v] {
			if !visited[u] {
				visited[u] = true
				parent[u] = v
				queue = append(queue, u)
			}
		}
	}

	canon := make([]string, n+1)
	count := 0

	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		childCanons := []string{}
		for _, u := range adj[v] {
			if u != parent[v] {
				childCanons = append(childCanons, canon[u])
			}
		}
		sort.Strings(childCanons)
		var sb strings.Builder
		sb.WriteByte('(')
		for _, c := range childCanons {
			sb.WriteString(c)
		}
		sb.WriteByte(')')
		canon[v] = sb.String()
		if inputTypes[canon[v]] {
			count++
		}
	}
	return count
}

// --- Infrastructure ---

func buildReference() (string, error) {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}

	srcBytes, err := os.ReadFile(refPath)
	if err != nil {
		return "", fmt.Errorf("cannot read reference source: %v", err)
	}
	srcContent := string(srcBytes)

	tmp, err := os.CreateTemp("", "1867F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	if strings.Contains(srcContent, "#include") {
		// C++ source saved with .go extension: copy to .cpp, compile with g++.
		cppPath := tmp.Name() + ".cpp"
		if err := os.WriteFile(cppPath, srcBytes, 0644); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("failed to write cpp source: %v", err)
		}
		defer os.Remove(cppPath)
		cmd := exec.Command("g++", "-O2", "-o", tmp.Name(), cppPath)
		var combined bytes.Buffer
		cmd.Stdout = &combined
		cmd.Stderr = &combined
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, combined.String())
		}
	} else {
		// Go source
		cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
		var combined bytes.Buffer
		cmd.Stdout = &combined
		cmd.Stderr = &combined
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, combined.String())
		}
	}
	return tmp.Name(), nil
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
