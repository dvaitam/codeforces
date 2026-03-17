package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct {
	u, v int
	w    int
}

type testCase struct {
	n     int
	edges []edge
}

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func (d *dsu) connected(n int) bool {
	root := d.find(1)
	for i := 2; i <= n; i++ {
		if d.find(i) != root {
			return false
		}
	}
	return true
}

func solve(tc testCase) string {
	cur := tc.edges
	ans := 0
	for bit := 29; bit >= 0; bit-- {
		d := newDSU(tc.n)
		nextEdges := make([]edge, 0, len(cur))
		mask := 1 << uint(bit)
		for _, e := range cur {
			if e.w&mask == 0 {
				d.union(e.u, e.v)
				nextEdges = append(nextEdges, e)
			}
		}
		if d.connected(tc.n) {
			cur = nextEdges
		} else {
			ans |= mask
		}
	}
	return fmt.Sprint(ans)
}

func generateTests(rng *rand.Rand) []testCase {
	var cases []testCase
	for len(cases) < 100 {
		n := rng.Intn(6) + 2 // 2..7 nodes
		// Generate a random spanning tree first to ensure connectivity
		edges := make([]edge, 0)
		for i := 2; i <= n; i++ {
			u := rng.Intn(i-1) + 1
			w := rng.Intn(1024) // weights up to 1023
			edges = append(edges, edge{u: u, v: i, w: w})
		}
		// Add some extra random edges
		extra := rng.Intn(n)
		for j := 0; j < extra; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			w := rng.Intn(1024)
			edges = append(edges, edge{u: u, v: v, w: w})
		}
		cases = append(cases, testCase{n: n, edges: edges})
	}
	return cases
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateTests(rng)

	for idx, tc := range cases {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
