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

// buildParent returns parent array for tree rooted at 1 (BFS).
func buildParent(adj [][]int, n int) []int {
	parent := make([]int, n+1)
	visited := make([]bool, n+1)
	visited[1] = true
	queue := []int{1}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, w := range adj[u] {
			if !visited[w] {
				visited[w] = true
				parent[w] = u
				queue = append(queue, w)
			}
		}
	}
	return parent
}

// oracle computes max c(v) - m*k by trying all subsets of deletable vertices.
// Deletable = all vertices except 1 (root) and v.
// c(v) after deleting set S = number of u not in S, u != v,
// whose first non-deleted ancestor (walking toward root) is v.
func oracle(n int, adj [][]int, v, k int) int {
	parent := buildParent(adj, n)

	deletable := []int{}
	for u := 1; u <= n; u++ {
		if u != 1 && u != v {
			deletable = append(deletable, u)
		}
	}

	computeCV := func(deleted []bool) int {
		count := 0
		for u := 1; u <= n; u++ {
			if u == v || deleted[u] {
				continue
			}
			cur := parent[u]
			for cur != 0 && deleted[cur] {
				cur = parent[cur]
			}
			if cur == v {
				count++
			}
		}
		return count
	}

	deleted := make([]bool, n+1)
	best := computeCV(deleted)
	nd := len(deletable)
	for mask := 1; mask < (1 << nd); mask++ {
		m := 0
		for i := 0; i < nd; i++ {
			if mask>>i&1 == 1 {
				deleted[deletable[i]] = true
				m++
			}
		}
		score := computeCV(deleted) - m*k
		if score > best {
			best = score
		}
		for i := 0; i < nd; i++ {
			deleted[deletable[i]] = false
		}
	}
	return best
}

func parseInput(input string) (int, [][]int, []int, []int) {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(r, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var q int
	fmt.Fscan(r, &q)
	vs := make([]int, q)
	ks := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(r, &vs[i], &ks[i])
	}
	return n, adj, vs, ks
}

func genCase(r *rand.Rand) string {
	n := r.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		parent := r.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", parent, i))
	}
	q := r.Intn(8) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		v := r.Intn(n) + 1
		k := r.Intn(6)
		sb.WriteString(fmt.Sprintf("%d %d\n", v, k))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 200; i++ {
		input := genCase(rng)
		n, adj, vs, ks := parseInput(input)

		var expParts []string
		for j := range vs {
			ans := oracle(n, adj, vs[j], ks[j])
			expParts = append(expParts, fmt.Sprintf("%d", ans))
		}
		expect := strings.Join(expParts, "\n")

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 200 tests passed")
}
