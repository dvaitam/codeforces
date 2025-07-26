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

func runCandidate(bin, input string) (string, error) {
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

func getPath(u, v int, adj [][]int) []int {
	if u == v {
		return []int{u}
	}
	n := len(adj) - 1
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = -1
	}
	q := []int{u}
	parent[u] = 0
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if x == v {
			break
		}
		for _, nb := range adj[x] {
			if parent[nb] == -1 {
				parent[nb] = x
				q = append(q, nb)
			}
		}
	}
	path := []int{}
	cur := v
	for cur != u {
		path = append(path, cur)
		cur = parent[cur]
	}
	path = append(path, u)
	return path
}

func sumDistinct(n int, adj [][]int, color []int) int64 {
	total := int64(0)
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			p := getPath(i, j, adj)
			seen := make(map[int]bool)
			for _, node := range p {
				seen[color[node]] = true
			}
			total += int64(len(seen))
		}
	}
	return total
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	m := rng.Intn(4) + 1
	adj := make([][]int, n+1)
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		edges = append(edges, [2]int{u, v})
	}
	color := make([]int, n+1)
	for i := 1; i <= n; i++ {
		color[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", color[i]))
	}
	sb.WriteByte('\n')
	res := make([]int64, 0, m+1)
	res = append(res, sumDistinct(n, adj, color))
	for k := 0; k < m; k++ {
		u := rng.Intn(n) + 1
		x := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", u, x))
		color[u] = x
		res = append(res, sumDistinct(n, adj, color))
	}
	var out strings.Builder
	for i, v := range res {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
