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

type edge struct{ u, v int }

type testCaseE struct {
	n     int
	edges []edge
	q     int
	ops   []opE
}

type opE struct {
	typ int // 0 or 1
	v   int
	x   int
	d   int
}

func generateTreeE(rng *rand.Rand, n int) []edge {
	deg := make([]int, n+1)
	var edges []edge
	for i := 2; i <= n; i++ {
		for {
			p := rng.Intn(i-1) + 1
			if p == 1 || deg[p] < 2 {
				edges = append(edges, edge{p, i})
				deg[p]++
				deg[i]++
				break
			}
		}
	}
	return edges
}

func generateCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(6) + 1
	edges := generateTreeE(rng, n)
	q := rng.Intn(10) + 1
	ops := make([]opE, q)
	for i := 0; i < q; i++ {
		typ := rng.Intn(2)
		if typ == 0 {
			v := rng.Intn(n) + 1
			x := rng.Intn(10)
			d := rng.Intn(n)
			ops[i] = opE{typ: 0, v: v, x: x, d: d}
		} else {
			v := rng.Intn(n) + 1
			ops[i] = opE{typ: 1, v: v}
		}
	}
	return testCaseE{n: n, edges: edges, q: q, ops: ops}
}

func buildInputE(tc testCaseE) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for _, op := range tc.ops {
		if op.typ == 0 {
			fmt.Fprintf(&sb, "0 %d %d %d\n", op.v, op.x, op.d)
		} else {
			fmt.Fprintf(&sb, "1 %d\n", op.v)
		}
	}
	return sb.String()
}

func expectedE(tc testCaseE) string {
	adj := make([][]int, tc.n+1)
	for _, e := range tc.edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	val := make([]int64, tc.n+1)
	var out strings.Builder
	for _, op := range tc.ops {
		if op.typ == 0 {
			// BFS
			seen := make([]bool, tc.n+1)
			queue := []struct{ v, d int }{{op.v, 0}}
			seen[op.v] = true
			for len(queue) > 0 {
				cur := queue[0]
				queue = queue[1:]
				val[cur.v] += int64(op.x)
				if cur.d == op.d {
					continue
				}
				for _, nb := range adj[cur.v] {
					if !seen[nb] {
						seen[nb] = true
						queue = append(queue, struct{ v, d int }{nb, cur.d + 1})
					}
				}
			}
		} else {
			fmt.Fprintf(&out, "%d\n", val[op.v])
		}
	}
	return strings.TrimSuffix(out.String(), "\n")
}

func runCaseE(bin string, tc testCaseE) error {
	input := buildInputE(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := expectedE(tc)
	if got != want {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInputE(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
