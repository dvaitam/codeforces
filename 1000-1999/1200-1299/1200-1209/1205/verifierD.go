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

	tmp, err := os.CreateTemp("", "1205D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	if strings.Contains(srcContent, "#include") {
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

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// parseInput parses the tree input and returns n and the edge list.
func parseInput(input string) (int, [][2]int) {
	fields := strings.Fields(input)
	n, _ := strconv.Atoi(fields[0])
	var edges [][2]int
	idx := 1
	for i := 0; i < n-1; i++ {
		u, _ := strconv.Atoi(fields[idx])
		v, _ := strconv.Atoi(fields[idx+1])
		edges = append(edges, [2]int{u, v})
		idx += 2
	}
	return n, edges
}

// validateOutput checks that the candidate output is a valid edge-weight assignment
// for the given tree, with |w| <= n*n and distinct path sums >= 2*n*n/9.
func validateOutput(n int, inputEdges [][2]int, output string) error {
	if n == 1 {
		// No edges, output should be empty
		if strings.TrimSpace(output) != "" {
			return fmt.Errorf("expected empty output for n=1, got: %s", output)
		}
		return nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != n-1 {
		return fmt.Errorf("expected %d lines, got %d", n-1, len(lines))
	}

	// Build adjacency from input edges (as a set of undirected edges)
	type edgeKey struct{ u, v int }
	inputEdgeSet := make(map[edgeKey]bool)
	for _, e := range inputEdges {
		u, v := e[0], e[1]
		if u > v {
			u, v = v, u
		}
		inputEdgeSet[edgeKey{u, v}] = true
	}

	// Parse output edges and build weighted adjacency list
	type wedge struct {
		to, w int
	}
	adj := make(map[int][]wedge)
	usedEdges := make(map[edgeKey]bool)
	maxW := n * n

	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return fmt.Errorf("line %d: expected 3 fields, got %d", i+1, len(fields))
		}
		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		w, err3 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("line %d: parse error", i+1)
		}
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("line %d: vertex out of range", i+1)
		}
		if w < -maxW || w > maxW {
			return fmt.Errorf("line %d: weight %d out of range [-%d, %d]", i+1, w, maxW, maxW)
		}
		// Check edge exists in input
		eu, ev := u, v
		if eu > ev {
			eu, ev = ev, eu
		}
		key := edgeKey{eu, ev}
		if !inputEdgeSet[key] {
			return fmt.Errorf("line %d: edge (%d, %d) not in input tree", i+1, u, v)
		}
		if usedEdges[key] {
			return fmt.Errorf("line %d: duplicate edge (%d, %d)", i+1, u, v)
		}
		usedEdges[key] = true
		adj[u] = append(adj[u], wedge{v, w})
		adj[v] = append(adj[v], wedge{u, w})
	}

	if len(usedEdges) != n-1 {
		return fmt.Errorf("expected %d edges, got %d", n-1, len(usedEdges))
	}

	// Compute all pairwise path sums using BFS/DFS from each node
	distinctSums := make(map[int]bool)
	for start := 1; start <= n; start++ {
		// BFS from start
		dist := make(map[int]int)
		dist[start] = 0
		queue := []int{start}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, e := range adj[u] {
				if _, visited := dist[e.to]; !visited {
					dist[e.to] = dist[u] + e.w
					queue = append(queue, e.to)
				}
			}
		}
		for v, d := range dist {
			if v > start {
				distinctSums[d] = true
			}
		}
	}

	count := len(distinctSums)
	// The problem requires at least floor(2*n*n/9) distinct sums
	// Using the same threshold as the reference solution
	threshold := 2 * n * n / 9
	if count < threshold {
		return fmt.Errorf("only %d distinct path sums, need at least %d (2n^2/9 where n=%d)", count, threshold, n)
	}
	return nil
}

func genTree(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 2
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i) + 1
		edges[i-1] = [2]int{p, i + 1}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// Also build and validate reference to ensure our test infra works
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 100
	for i := 1; i <= total; i++ {
		input, n := genTree(rng)
		_, inputEdges := parseInput(input)

		// Run reference to make sure it works (sanity check)
		refOut, err := run(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := validateOutput(n, inputEdges, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on case %d: %v\n", i, err)
			os.Exit(1)
		}

		// Run candidate and validate
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if err := validateOutput(n, inputEdges, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", total)
}
