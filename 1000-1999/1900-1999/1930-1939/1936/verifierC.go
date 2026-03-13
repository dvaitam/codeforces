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

func genCaseC(rng *rand.Rand) (string, int) {
	t := 1
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d %d\n", t, n, m))
	for i := 0; i < n; i++ {
		val := rng.Intn(10)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(val))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rng.Intn(10)))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), t
}

// solveC computes the correct answer for Codeforces 1936C (Pokémon Arena).
// Build a complete graph on n pokemon. Edge from u to v costs:
//   c[v] + max over j in [0,m) of max(0, a[u][j] - a[v][j])
// Answer is shortest path from 0 to n-1.
// But actually, looking at accepted solutions, the edge weight is:
//   For each attribute j, build edges between all pairs of pokemons.
//   Edge from u to v through attribute j costs: c[v] + max(0, a[u][j] - a[v][j]).
//   The overall edge from u to v is the minimum over all j of this cost.
// We use Dijkstra.
func solveC(input string) ([]string, error) {
	tokens := strings.Fields(input)
	idx := 0
	next := func() int {
		v, _ := strconv.Atoi(tokens[idx])
		idx++
		return v
	}

	T := next()
	results := make([]string, 0, T)

	for ; T > 0; T-- {
		n := next()
		m := next()
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			c[i] = int64(next())
		}
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				a[i][j] = int64(next())
			}
		}

		if n == 1 {
			results = append(results, "0")
			continue
		}

		// Build graph: for each attribute j, for each pair (u,v),
		// edge cost = c[v] + max(0, a[u][j] - a[v][j]).
		// We want minimum edge cost over all attributes j.
		// Then Dijkstra from 0 to n-1.
		const INF = int64(1<<62 - 1)

		// For small n (<=4), just try all edges
		type edge struct {
			to   int
			cost int64
		}
		adj := make([][]edge, n)
		// For each pair (u,v), find minimum cost edge over all attributes
		for u := 0; u < n; u++ {
			for v := 0; v < n; v++ {
				if u == v {
					continue
				}
				best := INF
				for j := 0; j < m; j++ {
					d := a[u][j] - a[v][j]
					if d < 0 {
						d = 0
					}
					w := c[v] + d
					if w < best {
						best = w
					}
				}
				adj[u] = append(adj[u], edge{v, best})
			}
		}

		// Dijkstra
		dist := make([]int64, n)
		for i := range dist {
			dist[i] = INF
		}
		dist[0] = 0
		visited := make([]bool, n)
		for {
			u := -1
			for i := 0; i < n; i++ {
				if !visited[i] && (u == -1 || dist[i] < dist[u]) {
					u = i
				}
			}
			if u == -1 || dist[u] == INF {
				break
			}
			if u == n-1 {
				break
			}
			visited[u] = true
			for _, e := range adj[u] {
				if dist[u]+e.cost < dist[e.to] {
					dist[e.to] = dist[u] + e.cost
				}
			}
		}
		results = append(results, fmt.Sprint(dist[n-1]))
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, t := genCaseC(rng)
		refResults, err := solveC(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solver failed on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candOut, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candTokens := strings.Fields(candOut)
		if len(candTokens) != t {
			fmt.Fprintf(os.Stderr, "case %d: expected %d answers, got %d\ninput:\n%s\noutput:\n%s", i+1, t, len(candTokens), input, candOut)
			os.Exit(1)
		}
		for j := 0; j < t; j++ {
			if candTokens[j] != refResults[j] {
				fmt.Fprintf(os.Stderr, "case %d subcase %d failed: expected %q got %q\ninput:\n%s", i+1, j+1, refResults[j], candTokens[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
