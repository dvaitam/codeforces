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

func runCandidate(bin, input string) (string, error) {
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

func buildBCC(n int, edges []edge) []int {
	adj := make([][]struct{ to, id int }, n)
	for i, e := range edges {
		u, v := e.u, e.v
		adj[u] = append(adj[u], struct{ to, id int }{v, i})
		adj[v] = append(adj[v], struct{ to, id int }{u, i})
	}
	timer := 0
	tin := make([]int, n)
	low := make([]int, n)
	bridge := make([]bool, len(edges))
	var dfs func(int, int)
	dfs = func(v, pe int) {
		timer++
		tin[v] = timer
		low[v] = timer
		for _, e := range adj[v] {
			if e.id == pe {
				continue
			}
			if tin[e.to] == 0 {
				dfs(e.to, e.id)
				if low[e.to] > tin[v] {
					bridge[e.id] = true
				}
				if low[e.to] < low[v] {
					low[v] = low[e.to]
				}
			} else if tin[e.to] < low[v] {
				low[v] = tin[e.to]
			}
		}
	}
	for i := 0; i < n; i++ {
		if tin[i] == 0 {
			dfs(i, -1)
		}
	}
	comp := make([]int, n)
	cid := 0
	var dfs2 func(int)
	dfs2 = func(v int) {
		comp[v] = cid
		for _, e := range adj[v] {
			if comp[e.to] == 0 && !bridge[e.id] {
				dfs2(e.to)
			}
		}
	}
	for i := 0; i < n; i++ {
		if comp[i] == 0 {
			cid++
			dfs2(i)
		}
	}
	return comp
}

func checkQuery(n int, base []edge, fav []int, qEdges []edge) bool {
	allEdges := append([]edge(nil), base...)
	allEdges = append(allEdges, qEdges...)
	comp := buildBCC(n, allEdges)
	c0 := comp[fav[0]]
	for _, v := range fav[1:] {
		if comp[v] != c0 {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(6)
	base := make([]edge, m)
	for i := 0; i < m; i++ {
		base[i] = edge{rng.Intn(n), rng.Intn(n)}
	}
	q := rng.Intn(3) + 1
	input := fmt.Sprintf("%d %d %d\n", n, m, q)
	for _, e := range base {
		input += fmt.Sprintf("%d %d\n", e.u+1, e.v+1)
	}
	answers := make([]string, q)
	for qi := 0; qi < q; qi++ {
		ni := rng.Intn(n) + 1
		mi := rng.Intn(3)
		favSet := rand.Perm(n)[:ni]
		input += fmt.Sprintf("%d %d\n", ni, mi)
		for i, idx := range favSet {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", idx+1)
		}
		input += "\n"
		qEdges := make([]edge, mi)
		for i := 0; i < mi; i++ {
			qEdges[i] = edge{rng.Intn(n), rng.Intn(n)}
			input += fmt.Sprintf("%d %d\n", qEdges[i].u+1, qEdges[i].v+1)
		}
		if checkQuery(n, base, favSet, qEdges) {
			answers[qi] = "YES"
		} else {
			answers[qi] = "NO"
		}
	}
	expected := strings.Join(answers, "\n")
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
