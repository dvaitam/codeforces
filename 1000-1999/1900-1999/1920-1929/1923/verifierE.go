package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct{ u, v int }

type test struct {
	n      int
	colors []int
	edges  []edge
}

func countPaths(n int, colors []int, edges []edge) int64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	var ans int64
	// BFS for each pair
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if colors[i-1] != colors[j-1] {
				continue
			}
			// BFS to find path
			prev := make([]int, n+1)
			for idx := range prev {
				prev[idx] = -1
			}
			q := []int{i}
			prev[i] = i
			for len(q) > 0 && prev[j] == -1 {
				cur := q[0]
				q = q[1:]
				for _, to := range adj[cur] {
					if prev[to] == -1 {
						prev[to] = cur
						q = append(q, to)
					}
				}
			}
			if prev[j] == -1 {
				continue
			}
			ok := true
			cur := j
			for prev[cur] != cur {
				cur = prev[cur]
				if cur != i && colors[cur-1] == colors[i-1] {
					ok = false
					break
				}
			}
			if ok {
				ans++
			}
		}
	}
	return ans
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(6) + 2
		colors := make([]int, n)
		for i := 0; i < n; i++ {
			colors[i] = rng.Intn(3) + 1
		}
		edges := make([]edge, 0, n-1)
		parent := make([]int, n+1)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges = append(edges, edge{p, i})
			parent[i] = p
		}
		tests = append(tests, test{n, colors, edges})
	}
	return tests
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(t.n))
		sb.WriteString("\n")
		for j, v := range t.colors {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteString("\n")
		for _, e := range t.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		expected := strconv.FormatInt(countPaths(t.n, t.colors, t.edges), 10)
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
