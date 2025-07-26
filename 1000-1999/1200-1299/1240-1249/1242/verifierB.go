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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type edge struct{ a, b int }

func solve(n int, edges []edge) int {
	adj := make([]map[int]struct{}, n+1)
	for _, e := range edges {
		if adj[e.a] == nil {
			adj[e.a] = make(map[int]struct{})
		}
		if adj[e.b] == nil {
			adj[e.b] = make(map[int]struct{})
		}
		adj[e.a][e.b] = struct{}{}
		adj[e.b][e.a] = struct{}{}
	}
	unvis := make(map[int]struct{}, n)
	for i := 1; i <= n; i++ {
		unvis[i] = struct{}{}
	}
	components := 0
	queue := make([]int, 0)
	for len(unvis) > 0 {
		var start int
		for k := range unvis {
			start = k
			break
		}
		queue = append(queue, start)
		delete(unvis, start)
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for u := range unvis {
				if _, ok := adj[v][u]; !ok {
					queue = append(queue, u)
					delete(unvis, u)
				}
			}
		}
		components++
	}
	return components - 1
}

func genCase(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(20) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edges := make([]edge, 0, m)
	used := make(map[[2]int]struct{})
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		edges = append(edges, edge{a, b})
	}
	return n, edges
}

func formatInput(n int, edges []edge) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.a, e.b))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct {
		n  int
		es []edge
	}
	// simple base cases
	cases = append(cases, struct {
		n  int
		es []edge
	}{1, nil})
	cases = append(cases, struct {
		n  int
		es []edge
	}{2, []edge{{1, 2}}})
	cases = append(cases, struct {
		n  int
		es []edge
	}{3, []edge{}})
	for len(cases) < 100 {
		n, es := genCase(rng)
		cases = append(cases, struct {
			n  int
			es []edge
		}{n, es})
	}

	for i, tc := range cases {
		in := formatInput(tc.n, tc.es)
		exp := fmt.Sprintf("%d", solve(tc.n, tc.es))
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
