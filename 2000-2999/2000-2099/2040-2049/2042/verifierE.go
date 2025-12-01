package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const refSource2042E = "2042E.go"

type testCase struct {
	n      int
	values []int
	edges  [][2]int
	adj    [][]int
}

type namedCase struct {
	name string
	tc   testCase
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource2042E)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, item := range tests {
		input := buildInput(item.tc)

		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, item.name, err, input)
			os.Exit(1)
		}
		refSet, err := parseOutput(refOut, item.tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error: could not parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, item.name, err, refOut)
			os.Exit(1)
		}
		if err := validateSubset(item.tc, refSet); err != nil {
			fmt.Fprintf(os.Stderr, "internal error: reference produced invalid subset on test %d (%s): %v\noutput:\n%s", idx+1, item.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, item.name, err, input)
			os.Exit(1)
		}
		candSet, err := parseOutput(candOut, item.tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: cannot parse output: %v\ninput:\n%soutput:\n%s", idx+1, item.name, err, input, candOut)
			os.Exit(1)
		}
		if err := validateSubset(item.tc, candSet); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: invalid subset: %v\ninput:\n%soutput:\n%s", idx+1, item.name, err, input, candOut)
			os.Exit(1)
		}

		if err := compareSets(refSet, candSet); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference subset:\n%v\ncandidate subset:\n%v\n", idx+1, item.name, err, input, refSet, candSet)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier2042E-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		source := filepath.Join(".", filepath.Clean(path))
		cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []namedCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []namedCase{
		{name: "tiny_path", tc: pathCase(3)},
		{name: "star", tc: starCase(5)},
		{name: "random_small", tc: randomCase(rng, 8)},
		{name: "random_medium", tc: randomCase(rng, 120)},
		{name: "heavy", tc: randomCase(rng, 60000)},
	}
}

func pathCase(n int) testCase {
	// vertices = 2n, connect linearly
	total := 2 * n
	edges := make([][2]int, 0, total-1)
	for i := 2; i <= total; i++ {
		edges = append(edges, [2]int{i - 1, i})
	}
	vals := pairedValues(n)
	return testCase{n: n, values: vals, edges: edges, adj: buildAdj(total, edges)}
}

func starCase(n int) testCase {
	total := 2 * n
	edges := make([][2]int, 0, total-1)
	for i := 2; i <= total; i++ {
		edges = append(edges, [2]int{1, i})
	}
	vals := pairedValues(n)
	return testCase{n: n, values: vals, edges: edges, adj: buildAdj(total, edges)}
}

func randomCase(rng *rand.Rand, n int) testCase {
	total := 2 * n
	edges := make([][2]int, 0, total-1)
	for v := 2; v <= total; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	vals := pairedValues(n)
	rng.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	return testCase{n: n, values: vals, edges: edges, adj: buildAdj(total, edges)}
}

func pairedValues(n int) []int {
	vals := make([]int, 2*n)
	for i := 0; i < n; i++ {
		vals[2*i] = i + 1
		vals[2*i+1] = i + 1
	}
	return vals
}

func buildAdj(total int, edges [][2]int) [][]int {
	adj := make([][]int, total+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	return adj
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d", e[0], e[1])
		if i+1 != len(tc.edges) {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string, tc testCase) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid k: %v", err)
	}
	if k < 0 || k > 2*tc.n {
		return nil, fmt.Errorf("k out of range: %d", k)
	}
	if len(fields) != k+1 {
		return nil, fmt.Errorf("expected %d vertex ids, got %d", k, len(fields)-1)
	}
	sub := make([]int, k)
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid vertex id %q", fields[i+1])
		}
		sub[i] = v
	}
	return sub, nil
}

func validateSubset(tc testCase, subset []int) error {
	total := 2 * tc.n
	if len(subset) == 0 {
		return fmt.Errorf("subset is empty")
	}
	seen := make([]bool, total+1)
	for _, v := range subset {
		if v < 1 || v > total {
			return fmt.Errorf("vertex %d out of range", v)
		}
		if seen[v] {
			return fmt.Errorf("vertex %d duplicated", v)
		}
		seen[v] = true
	}

	covered := make([]bool, tc.n+1)
	for v := range tc.values {
		if seen[v+1] {
			covered[tc.values[v]] = true
		}
	}
	for i := 1; i <= tc.n; i++ {
		if !covered[i] {
			return fmt.Errorf("value %d missing in subset", i)
		}
	}

	// connectivity check with BFS inside subset
	queue := make([]int, 0, len(subset))
	queue = append(queue, subset[0])
	vis := make([]bool, total+1)
	vis[subset[0]] = true
	for front := 0; front < len(queue); front++ {
		u := queue[front]
		for _, nxt := range tc.adj[u] {
			if !seen[nxt] || vis[nxt] {
				continue
			}
			vis[nxt] = true
			queue = append(queue, nxt)
		}
	}
	if len(queue) != len(subset) {
		return fmt.Errorf("subset not connected")
	}

	return nil
}

func compareSets(expected, actual []int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("subset size mismatch: expected %d got %d", len(expected), len(actual))
	}
	exp := append([]int(nil), expected...)
	act := append([]int(nil), actual...)
	sort.Ints(exp)
	sort.Ints(act)
	for i := range exp {
		if exp[i] != act[i] {
			return fmt.Errorf("subset differs at position %d: expected %d got %d", i, exp[i], act[i])
		}
	}
	return nil
}
