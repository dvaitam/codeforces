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
)

func buildReference() (string, error) {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	outBin := filepath.Join(os.TempDir(), "ref553D")
	content, err := os.ReadFile(refPath)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	if strings.Contains(string(content), "#include") {
		cppPath := filepath.Join(os.TempDir(), "ref553D.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("g++", "-O2", "-o", outBin, cppPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build ref (c++) failed: %v\n%s", err, o)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", outBin, refPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build ref failed: %v\n%s", err, o)
		}
	}
	return outBin, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type graph struct {
	n, m, k int
	fortress map[int]bool
	adj      [][]int
	input    string
}

func generateCase(rng *rand.Rand) graph {
	n := rng.Intn(8) + 3 // 3..10
	k := rng.Intn(n-1) + 1 // 1..n-1
	// pick k fortress cities
	perm := rng.Perm(n)
	fortress := make(map[int]bool)
	for i := 0; i < k; i++ {
		fortress[perm[i]+1] = true
	}
	// generate edges ensuring each city has at least one neighbor
	type edge struct{ u, v int }
	edgeSet := make(map[edge]bool)
	// first ensure connectivity-ish: each node gets at least one edge
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		e := edge{u, i}
		if u > i {
			e = edge{i, u}
		}
		edgeSet[e] = true
	}
	// add some random edges
	extra := rng.Intn(n)
	for i := 0; i < extra; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		edgeSet[edge{u, v}] = true
	}
	m := len(edgeSet)
	adj := make([][]int, n+1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	flist := make([]int, 0, k)
	for c := range fortress {
		flist = append(flist, c)
	}
	for i, c := range flist {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", c)
	}
	sb.WriteByte('\n')
	for e := range edgeSet {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	return graph{n: n, m: m, k: k, fortress: fortress, adj: adj, input: sb.String()}
}

func validateOutput(g graph, output string) (float64, error) {
	lines := strings.Fields(output)
	if len(lines) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	r, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, fmt.Errorf("invalid r: %v", err)
	}
	if r < 1 || r > g.n-g.k {
		return 0, fmt.Errorf("invalid r=%d (n=%d, k=%d)", r, g.n, g.k)
	}
	if len(lines) != r+1 {
		return 0, fmt.Errorf("expected %d cities, got %d tokens", r, len(lines)-1)
	}
	inS := make(map[int]bool)
	for i := 1; i <= r; i++ {
		c, err := strconv.Atoi(lines[i])
		if err != nil {
			return 0, fmt.Errorf("invalid city: %v", err)
		}
		if c < 1 || c > g.n {
			return 0, fmt.Errorf("city %d out of range", c)
		}
		if g.fortress[c] {
			return 0, fmt.Errorf("city %d is a fortress", c)
		}
		if inS[c] {
			return 0, fmt.Errorf("duplicate city %d", c)
		}
		inS[c] = true
	}
	minStrength := math.Inf(1)
	for c := range inS {
		total := len(g.adj[c])
		if total == 0 {
			continue
		}
		inCount := 0
		for _, nb := range g.adj[c] {
			if inS[nb] {
				inCount++
			}
		}
		strength := float64(inCount) / float64(total)
		if strength < minStrength {
			minStrength = strength
		}
	}
	return minStrength, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(42))
	for idx := 0; idx < 50; idx++ {
		g := generateCase(rng)

		refOut, err := run(refBin, g.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		refStrength, err := validateOutput(g, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		got, err := run(bin, g.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotStrength, err := validateOutput(g, got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output: %v\ninput:\n%s", idx+1, err, g.input)
			os.Exit(1)
		}
		if gotStrength < refStrength-1e-9 {
			fmt.Fprintf(os.Stderr, "case %d failed: suboptimal strength %.6f < %.6f\ninput:\n%sref output:\n%s\ngot output:\n%s\n",
				idx+1, gotStrength, refStrength, g.input, refOut, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All 50 tests passed\n")
}
