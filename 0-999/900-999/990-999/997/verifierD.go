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

const mod int64 = 998244353

type testCase struct {
	input    string
	expected int64
}

type Edge struct {
	to  int
	rev int
}

type graph [][]Edge

func polyGeom(f []int64, k int) []int64 {
	res := make([]int64, k+1)
	res[0] = 1
	for d := 1; d <= k; d++ {
		var sum int64
		for i := 1; i <= d; i++ {
			sum += f[i] * res[d-i]
			sum %= mod
		}
		res[d] = sum
	}
	return res
}

func computeClosedWalks(n, k int, edges [][2]int) []int64 {
	adj := make(graph, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		ai := len(adj[a])
		bi := len(adj[b])
		adj[a] = append(adj[a], Edge{to: b, rev: bi})
		adj[b] = append(adj[b], Edge{to: a, rev: ai})
	}
	F := make([][][]int64, n)
	for v := 0; v < n; v++ {
		F[v] = make([][]int64, len(adj[v]))
		for i := range F[v] {
			F[v][i] = make([]int64, k+1)
		}
	}
	parent := make([]int, n)
	parentIdx := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	root := 0
	stack := []int{root}
	order := []int{}
	parent[root] = -1
	parentIdx[root] = -1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, e := range adj[v] {
			if e.to == parent[v] {
				continue
			}
			parent[e.to] = v
			parentIdx[e.to] = e.rev
			stack = append(stack, e.to)
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if parent[v] == -1 {
			continue
		}
		sum := make([]int64, k+1)
		for _, e := range adj[v] {
			if e.to == parent[v] {
				continue
			}
			for t := 0; t+2 <= k; t++ {
				sum[t+2] = (sum[t+2] + F[e.to][e.rev][t]) % mod
			}
		}
		F[v][parentIdx[v]] = polyGeom(sum, k)
	}
	cw := make([]int64, k+1)
	var dfs func(v, p int)
	dfs = func(v, p int) {
		deg := len(adj[v])
		prefix := make([][]int64, deg+1)
		suffix := make([][]int64, deg+1)
		for i := 0; i <= deg; i++ {
			prefix[i] = make([]int64, k+1)
			suffix[i] = make([]int64, k+1)
		}
		for i := 0; i < deg; i++ {
			copy(prefix[i+1], prefix[i])
			for t := 0; t+2 <= k; t++ {
				prefix[i+1][t+2] = (prefix[i+1][t+2] + F[adj[v][i].to][adj[v][i].rev][t]) % mod
			}
		}
		for i := deg - 1; i >= 0; i-- {
			copy(suffix[i], suffix[i+1])
			for t := 0; t+2 <= k; t++ {
				suffix[i][t+2] = (suffix[i][t+2] + F[adj[v][i].to][adj[v][i].rev][t]) % mod
			}
		}
		total := prefix[deg]
		P := polyGeom(total, k)
		for t := 0; t <= k; t++ {
			cw[t] = (cw[t] + P[t]) % mod
		}
		for idx := range adj[v] {
			sum := make([]int64, k+1)
			for t := 0; t <= k; t++ {
				sum[t] = (prefix[idx][t] + suffix[idx+1][t]) % mod
			}
			F[v][idx] = polyGeom(sum, k)
		}
		for _, e := range adj[v] {
			if e.to == p {
				continue
			}
			dfs(e.to, v)
		}
	}
	dfs(root, -1)
	return cw
}

func solveD(n1, n2, k int, edges1, edges2 [][2]int) int64 {
	cw1 := computeClosedWalks(n1, k, edges1)
	cw2 := computeClosedWalks(n2, k, edges2)
	C := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		C[i] = make([]int64, i+1)
		C[i][0] = 1
		C[i][i] = 1
		for j := 1; j < i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % mod
		}
	}
	var ans int64
	for i := 0; i <= k; i++ {
		ans = (ans + C[k][i]*cw1[i]%mod*cw2[k-i]) % mod
	}
	return ans
}

func genTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges[i-1] = [2]int{p, i}
	}
	return edges
}

func genCase(rng *rand.Rand) testCase {
	n1 := rng.Intn(4) + 1
	n2 := rng.Intn(4) + 1
	k := rng.Intn(4) + 1
	edges1 := genTree(rng, n1)
	edges2 := genTree(rng, n2)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n1, n2, k)
	for _, e := range edges1 {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	for _, e := range edges2 {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	expected := solveD(n1, n2, k, edges1, edges2)
	return testCase{input: sb.String(), expected: expected}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
