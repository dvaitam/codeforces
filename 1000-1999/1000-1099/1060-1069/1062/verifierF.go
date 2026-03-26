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

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

// Brute-force solver for small n: compute transitive closure, then check
// important / semi-important conditions.
func solveOracle(n, m int, edges []edge) string {
	// reach[u] is set of nodes reachable from u (including u itself)
	reach := make([][]bool, n+1)
	for u := 0; u <= n; u++ {
		reach[u] = make([]bool, n+1)
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
	}
	// BFS from each node
	for u := 1; u <= n; u++ {
		visited := reach[u]
		visited[u] = true
		queue := []int{u}
		for len(queue) > 0 {
			x := queue[0]
			queue = queue[1:]
			for _, v := range adj[x] {
				if !visited[v] {
					visited[v] = true
					queue = append(queue, v)
				}
			}
		}
	}

	// For each u, compute the "bad set" = nodes v != u such that
	// neither u->v nor v->u
	ans := 0
	for u := 1; u <= n; u++ {
		var badSet []int
		for v := 1; v <= n; v++ {
			if v == u {
				continue
			}
			if !reach[u][v] && !reach[v][u] {
				badSet = append(badSet, v)
			}
		}
		if len(badSet) == 0 {
			// important
			ans++
			continue
		}
		// semi-important: exists exactly one node w to destroy such that
		// u becomes important in the remaining graph
		// Equivalently: exists w != u such that for all v != u, v != w:
		// reach(u,v) or reach(v,u) in the original graph (removing w
		// can only help if w was blocking, but in a DAG removing w
		// removes edges through w).
		// Actually we need to check reachability in the graph with w removed.
		// For small n, just brute force it.
		semi := false
		for w := 1; w <= n; w++ {
			if w == u {
				continue
			}
			// Check if u is important in graph with w removed
			// Recompute reachability from u and to u without w
			reachFromU := make([]bool, n+1)
			reachFromU[u] = true
			q := []int{u}
			for len(q) > 0 {
				x := q[0]
				q = q[1:]
				for _, v := range adj[x] {
					if v != w && !reachFromU[v] {
						reachFromU[v] = true
						q = append(q, v)
					}
				}
			}
			reachToU := make([]bool, n+1)
			reachToU[u] = true
			// Build reverse adj
			q = []int{u}
			for len(q) > 0 {
				x := q[0]
				q = q[1:]
				for src := 1; src <= n; src++ {
					if src == w || reachToU[src] {
						continue
					}
					for _, v := range adj[src] {
						if v == x {
							reachToU[src] = true
							q = append(q, src)
							break
						}
					}
				}
			}
			// Actually let me do reverse adj properly for efficiency
			// For small n this is fine but let me redo with radj
			radj := make([][]int, n+1)
			for _, e := range edges {
				if e.u != w && e.v != w {
					radj[e.v] = append(radj[e.v], e.u)
				}
			}
			reachToU2 := make([]bool, n+1)
			reachToU2[u] = true
			q2 := []int{u}
			for len(q2) > 0 {
				x := q2[0]
				q2 = q2[1:]
				for _, v := range radj[x] {
					if !reachToU2[v] {
						reachToU2[v] = true
						q2 = append(q2, v)
					}
				}
			}

			ok := true
			for v := 1; v <= n; v++ {
				if v == u || v == w {
					continue
				}
				if !reachFromU[v] && !reachToU2[v] {
					ok = false
					break
				}
			}
			if ok {
				semi = true
				break
			}
		}
		if semi {
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, int, int, []edge) {
	n := rng.Intn(8) + 2
	maxM := n * (n - 1) / 2
	m := rng.Intn(maxM + 1)
	var allEdges []edge
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			allEdges = append(allEdges, edge{i, j})
		}
	}
	rng.Shuffle(len(allEdges), func(i, j int) { allEdges[i], allEdges[j] = allEdges[j], allEdges[i] })
	allEdges = allEdges[:m]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range allEdges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String(), n, m, allEdges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []struct {
		input string
		n, m  int
		edges []edge
	}{
		{"2 0\n", 2, 0, nil},
		{"2 1\n1 2\n", 2, 1, []edge{{1, 2}}},
		{"3 2\n1 2\n2 3\n", 3, 2, []edge{{1, 2}, {2, 3}}},
	}
	for i := 0; i < 100; i++ {
		input, n, m, edges := genCase(rng)
		tests = append(tests, struct {
			input string
			n, m  int
			edges []edge
		}{input, n, m, edges})
	}

	for idx, t := range tests {
		exp := solveOracle(t.n, t.m, t.edges)
		got, err := runProg(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, t.input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\n got: %s\ninput:\n%s", idx+1, exp, got, t.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
