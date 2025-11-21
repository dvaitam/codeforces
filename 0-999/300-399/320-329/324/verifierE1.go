package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type edgeRecord struct {
	from   int
	to     int
	vision []int
}

type graphBuilder struct {
	n     int
	adj   []map[int]*edgeRecord
	edges []*edgeRecord
}

type graphData struct {
	n   int
	adj []map[int][]int
}

type testCase struct {
	name       string
	input      string
	graph      *graphData
	expectPath bool
	witness    []int
}

func newGraphBuilder(n int) *graphBuilder {
	adj := make([]map[int]*edgeRecord, n+1)
	for i := range adj {
		adj[i] = make(map[int]*edgeRecord)
	}
	return &graphBuilder{n: n, adj: adj}
}

func (g *graphBuilder) canAdd(u, v int) bool {
	if u == v {
		return false
	}
	if g.adj[u][v] != nil || g.adj[v][u] != nil {
		return false
	}
	return true
}

func (g *graphBuilder) addEdge(u, v int, vision []int) (*edgeRecord, bool) {
	if !g.canAdd(u, v) {
		return nil, false
	}
	rec := &edgeRecord{from: u, to: v}
	if vision != nil {
		rec.vision = append([]int(nil), vision...)
	}
	g.adj[u][v] = rec
	g.edges = append(g.edges, rec)
	return rec, true
}

func (g *graphBuilder) toGraphData() *graphData {
	adj := make([]map[int][]int, g.n+1)
	for i := range adj {
		adj[i] = make(map[int][]int)
	}
	for _, rec := range g.edges {
		seq := append([]int(nil), rec.vision...)
		adj[rec.from][rec.to] = seq
	}
	return &graphData{n: g.n, adj: adj}
}

func (g *graphBuilder) toInput() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", g.n, len(g.edges))
	for _, rec := range g.edges {
		seq := rec.vision
		fmt.Fprintf(&sb, "%d %d %d", rec.from, rec.to, len(seq))
		for _, val := range seq {
			fmt.Fprintf(&sb, " %d", val)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func sampleCase() testCase {
	builder := newGraphBuilder(6)
	builder.addEdge(1, 2, []int{1, 2})
	builder.addEdge(2, 3, []int{3})
	builder.addEdge(3, 4, []int{4, 5})
	builder.addEdge(4, 5, []int{})
	builder.addEdge(5, 3, []int{3})
	builder.addEdge(6, 1, []int{6})
	return testCase{
		name:       "sample",
		input:      builder.toInput(),
		graph:      builder.toGraphData(),
		expectPath: true,
		witness:    []int{6, 1, 2, 3},
	}
}

func emptyGraphCase(n int, name string) testCase {
	builder := newGraphBuilder(n)
	return testCase{
		name:       name,
		input:      builder.toInput(),
		graph:      builder.toGraphData(),
		expectPath: false,
	}
}

func impossibleSingleEdgeCase() testCase {
	builder := newGraphBuilder(2)
	builder.addEdge(1, 2, []int{2})
	return testCase{
		name:       "single_bad_edge",
		input:      builder.toInput(),
		graph:      builder.toGraphData(),
		expectPath: false,
	}
}

func partitionWithLimit(total, segments, limit int, rng *rand.Rand) []int {
	parts := make([]int, segments)
	remaining := total
	for i := 0; i < segments; i++ {
		leftSegments := segments - i
		maxFuture := limit * (leftSegments - 1)
		minVal := remaining - maxFuture
		if minVal < 0 {
			minVal = 0
		}
		maxVal := remaining
		if maxVal > limit {
			maxVal = limit
		}
		if maxVal < minVal {
			minVal = maxVal
		}
		val := minVal
		if maxVal > minVal {
			val += rng.Intn(maxVal - minVal + 1)
		}
		parts[i] = val
		remaining -= val
	}
	return parts
}

func addRandomEdges(builder *graphBuilder, rng *rand.Rand, count int) {
	attempts := 0
	for added := 0; added < count && attempts < count*8; attempts++ {
		u := rng.Intn(builder.n) + 1
		v := rng.Intn(builder.n) + 1
		if !builder.canAdd(u, v) {
			continue
		}
		lenVision := rng.Intn(builder.n + 1)
		seq := make([]int, lenVision)
		for i := 0; i < lenVision; i++ {
			seq[i] = rng.Intn(builder.n) + 1
		}
		builder.addEdge(u, v, seq)
		added++
	}
}

func randomPositiveCase(rng *rand.Rand, name string, forcedN int) testCase {
	for attempt := 0; attempt < 200; attempt++ {
		n := forcedN
		if n < 3 {
			n = rng.Intn(18) + 3
		}
		builder := newGraphBuilder(n)
		maxLen := 2 * n
		if maxLen < 2 {
			maxLen = 2
		}
		pathLen := 2
		if maxLen > 2 {
			pathLen += rng.Intn(maxLen - 1)
		}
		path := make([]int, pathLen)
		path[0] = rng.Intn(n) + 1
		pathEdges := make([]*edgeRecord, pathLen-1)
		valid := true
		for i := 1; i < pathLen; i++ {
			success := false
			for tries := 0; tries < 400; tries++ {
				next := rng.Intn(n) + 1
				if next == path[i-1] {
					continue
				}
				if !builder.canAdd(path[i-1], next) {
					continue
				}
				rec, _ := builder.addEdge(path[i-1], next, nil)
				pathEdges[i-1] = rec
				path[i] = next
				success = true
				break
			}
			if !success {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		segments := partitionWithLimit(len(path), len(pathEdges), n, rng)
		cur := 0
		for i, rec := range pathEdges {
			segLen := segments[i]
			rec.vision = append([]int(nil), path[cur:cur+segLen]...)
			cur += segLen
		}
		if cur != len(path) {
			continue
		}
		extra := rng.Intn(n)
		addRandomEdges(builder, rng, extra)
		return testCase{
			name:       name,
			input:      builder.toInput(),
			graph:      builder.toGraphData(),
			expectPath: true,
			witness:    append([]int(nil), path...),
		}
	}
	panic("failed to build random positive case")
}

func generateTests() []testCase {
	tests := []testCase{
		sampleCase(),
		emptyGraphCase(3, "empty_n3"),
		impossibleSingleEdgeCase(),
	}
	deterministicSeeds := []int64{1, 2, 3, 4, 5}
	for _, seed := range deterministicSeeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomPositiveCase(rng, fmt.Sprintf("deterministic_%d", seed), 0))
	}
	tests = append(tests, randomPositiveCase(rand.New(rand.NewSource(42)), "big_n", 50))
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomPositiveCase(rng, fmt.Sprintf("random_%d", len(tests)), 0))
		if len(tests)%10 == 0 {
			size := rng.Intn(4) + 2
			tests = append(tests, emptyGraphCase(size, fmt.Sprintf("empty_%d", len(tests))))
		}
	}
	return tests
}

func runProgram(target, input string) (string, error) {
	if !filepath.IsAbs(target) {
		if abs, err := filepath.Abs(target); err == nil {
			target = abs
		}
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validateOutput(tc testCase, out string) error {
	reader := strings.NewReader(out)
	var first int
	if _, err := fmt.Fscan(reader, &first); err != nil {
		return fmt.Errorf("failed to read first integer: %v\noutput:\n%s", err, out)
	}
	if first == 0 {
		if tc.expectPath {
			return fmt.Errorf("expected a valid path but got 0")
		}
		return nil
	}
	if !tc.expectPath {
		return fmt.Errorf("expected output 0 (no valid path) but got %d", first)
	}
	k := first
	if k < 2 {
		return fmt.Errorf("path must contain at least 2 shops, got %d", k)
	}
	if k > 2*tc.graph.n {
		return fmt.Errorf("path length %d exceeds limit 2n=%d", k, 2*tc.graph.n)
	}
	shops := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &shops[i]); err != nil {
			return fmt.Errorf("failed to read shop %d: %v", i+1, err)
		}
		if shops[i] < 1 || shops[i] > tc.graph.n {
			return fmt.Errorf("shop %d = %d is out of range", i+1, shops[i])
		}
	}
	var vision []int
	for i := 0; i < k-1; i++ {
		from := shops[i]
		to := shops[i+1]
		seq, ok := tc.graph.adj[from][to]
		if !ok {
			return fmt.Errorf("edge %d (%d -> %d) does not exist", i+1, from, to)
		}
		vision = append(vision, seq...)
	}
	if len(vision) != k {
		return fmt.Errorf("vision length %d does not match path length %d", len(vision), k)
	}
	for i := 0; i < k; i++ {
		if vision[i] != shops[i] {
			return fmt.Errorf("vision[%d]=%d differs from shop %d=%d", i+1, vision[i], i+1, shops[i])
		}
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]
	tests := generateTests()
	for i, tc := range tests {
		out, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) runtime error: %v\noutput:\n%s\ninput:\n%s", i+1, tc.name, err, out, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			if len(tc.witness) > 0 {
				fmt.Fprintf(os.Stderr, "example valid path (%d shops): %v\n", len(tc.witness), tc.witness)
			}
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
