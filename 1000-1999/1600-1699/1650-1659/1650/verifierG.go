package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildExecutable(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

func oracle(exe string, test string) (string, error) {
	return run(exe, "1\n"+test)
}

func genConnectedGraph(rng *rand.Rand, n, m int) [][2]int {
	edges := make([][2]int, 0, m)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		pa, pb := find(a), find(b)
		if pa != pb {
			parent[pb] = pa
		}
	}
	// create spanning tree first
	for i := 1; i < n; i++ {
		u := rng.Intn(i)
		edges = append(edges, [2]int{u, i})
		union(u, i)
	}
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		ok := true
		for _, e := range edges {
			if e[0] == u && e[1] == v {
				ok = false
				break
			}
		}
		if ok {
			edges = append(edges, [2]int{u, v})
		}
	}
	return edges
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	extra := rng.Intn(maxEdges - (n - 1) + 1)
	m := (n - 1) + extra
	s := rng.Intn(n) + 1
	t := rng.Intn(n) + 1
	for t == s {
		t = rng.Intn(n) + 1
	}
	edges := genConnectedGraph(rng, n, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()
	oracleExe, oracleCleanup, err := buildExecutable("1650G.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer oracleCleanup()
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		expected, err := oracle(oracleExe, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, "1\n"+tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
