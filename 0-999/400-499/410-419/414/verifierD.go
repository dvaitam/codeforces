package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func can(S int64, depths []int, prefix []int64, p int64) bool {
	if S == 0 {
		return true
	}
	n := int64(len(depths))
	if S > n {
		return false
	}
	minCost := int64(1<<62 - 1)
	for i := S; i <= n; i++ {
		sumSeg := prefix[i] - prefix[i-S]
		T := int64(depths[i-1])
		cost := S*T - sumSeg
		if cost < minCost {
			minCost = cost
		}
		if minCost <= p {
			return true
		}
	}
	return minCost <= p
}

func solveD(m int, k, p int64, edges [][2]int) int64 {
	adj := make([][]int, m+1)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	depth := make([]int, m+1)
	q := make([]int, 0, m)
	q = append(q, 1)
	parent := make([]int, m+1)
	parent[1] = -1
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			q = append(q, v)
		}
	}
	n := m - 1
	depths := make([]int, 0, n)
	for i := 2; i <= m; i++ {
		depths = append(depths, depth[i])
	}
	sort.Ints(depths)
	prefix := make([]int64, n+1)
	for i, d := range depths {
		prefix[i+1] = prefix[i] + int64(d)
	}
	if k > int64(n) {
		k = int64(n)
	}
	lo, hi := int64(0), k
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if can(mid, depths, prefix, p) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func generateCase(rng *rand.Rand) (string, int64) {
	m := rng.Intn(8) + 2
	edges := make([][2]int, 0, m-1)
	for v := 2; v <= m; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	k := int64(rng.Intn(m))
	p := int64(rng.Intn(10))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", m, k, p))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	ans := solveD(m, k, p, edges)
	return sb.String(), ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output\ninput:\n%soutput:\n%s", i+1, input, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
