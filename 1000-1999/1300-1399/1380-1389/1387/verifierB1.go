package main

import (
	"bufio"
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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB1")
	cmd := exec.Command("go", "build", "-o", oracle, "1387B1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(9) + 2 // at least 2 nodes
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		fmt.Fprintf(&sb, "%d %d\n", p, i)
	}
	return sb.String()
}

type edge struct {
	u int
	v int
}

func parseInput(input string) (int, []edge, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, nil, fmt.Errorf("failed to read n: %w", err)
	}
	edges := make([]edge, 0, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return 0, nil, fmt.Errorf("failed to read edge %d: %w", i+1, err)
		}
		edges = append(edges, edge{u: u, v: v})
	}
	return n, edges, nil
}

func parseOutput(out string, n int) (int, []int, error) {
	fields := strings.Fields(out)
	if len(fields) != n+1 {
		return 0, nil, fmt.Errorf("expected %d integers in output, got %d", n+1, len(fields))
	}
	total, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse total distance: %w", err)
	}
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return 0, nil, fmt.Errorf("failed to parse permutation at position %d: %w", i, err)
		}
		p[i] = v
	}
	return total, p, nil
}

func validatePermutation(p []int, n int) error {
	seen := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if p[i] < 1 || p[i] > n {
			return fmt.Errorf("value %d at position %d is out of range [1, %d]", p[i], i, n)
		}
		if seen[p[i]] {
			return fmt.Errorf("value %d appears multiple times", p[i])
		}
		if p[i] == i {
			return fmt.Errorf("villager %d stays in the same house", i)
		}
		seen[p[i]] = true
	}
	return nil
}

func sumDistances(n int, edges []edge, p []int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}

	const maxLog = 18 // enough for n <= 1e5
	up := make([][]int, maxLog)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	depth := make([]int, n+1)

	queue := []int{1}
	parent := make([]int, n+1)
	parent[1] = 1
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		up[0][u] = parent[u]
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			queue = append(queue, v)
		}
	}

	for j := 1; j < maxLog; j++ {
		for v := 1; v <= n; v++ {
			up[j][v] = up[j-1][up[j-1][v]]
		}
	}

	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for j := 0; j < maxLog; j++ {
			if (diff>>j)&1 == 1 {
				a = up[j][a]
			}
		}
		if a == b {
			return a
		}
		for j := maxLog - 1; j >= 0; j-- {
			if up[j][a] != up[j][b] {
				a = up[j][a]
				b = up[j][b]
			}
		}
		return up[0][a]
	}

	total := 0
	for i := 1; i <= n; i++ {
		x := lca(i, p[i])
		total += depth[i] + depth[p[i]] - 2*depth[x]
	}
	return total
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input := genCase(rng)
		n, edges, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad generated input on case %d: %v\n", i, err)
			os.Exit(1)
		}

		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		expectTotal, expectPerm, err := parseOutput(expect, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on case %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := validatePermutation(expectPerm, n); err != nil {
			fmt.Fprintf(os.Stderr, "oracle permutation invalid on case %d: %v\n", i, err)
			os.Exit(1)
		}
		if gotSum := sumDistances(n, edges, expectPerm); gotSum != expectTotal {
			fmt.Fprintf(os.Stderr, "oracle reported wrong total on case %d: reported %d computed %d\n", i, expectTotal, gotSum)
			os.Exit(1)
		}

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		gotTotal, gotPerm, err := parseOutput(got, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if err := validatePermutation(gotPerm, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid permutation: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		computed := sumDistances(n, edges, gotPerm)
		if computed != gotTotal {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong reported total, reported %d computed %d\ninput:\n%s", i, gotTotal, computed, input)
			os.Exit(1)
		}
		if gotTotal != expectTotal {
			fmt.Fprintf(os.Stderr, "case %d failed: expected minimum %d got %d\ninput:\n%s", i, expectTotal, gotTotal, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
