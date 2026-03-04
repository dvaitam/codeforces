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

type edge struct {
	u int
	v int
	w int64
}

type testCase struct {
	n     int
	edges []edge
	s     int
}

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func gain(w int64) int64 {
	if w <= 0 {
		return 0
	}
	// Max k such that k*(k-1)/2 <= w.
	lo, hi := int64(1), int64(1)
	for hi*(hi-1)/2 <= w {
		hi <<= 1
	}
	for lo+1 < hi {
		mid := (lo + hi) >> 1
		if mid*(mid-1)/2 <= w {
			lo = mid
		} else {
			hi = mid
		}
	}
	k := lo
	return k*w - (k*(k-1)*(k+1))/6
}

func solve(tc testCase) int64 {
	n := tc.n
	g := make([][]edge, n+1)
	rg := make([][]int, n+1)
	for _, e := range tc.edges {
		g[e.u] = append(g[e.u], e)
		rg[e.v] = append(rg[e.v], e.u)
	}

	// Kosaraju: first pass order.
	vis := make([]bool, n+1)
	order := make([]int, 0, n)
	var dfs1 func(int)
	dfs1 = func(u int) {
		vis[u] = true
		for _, e := range g[u] {
			if !vis[e.v] {
				dfs1(e.v)
			}
		}
		order = append(order, u)
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			dfs1(i)
		}
	}

	comp := make([]int, n+1)
	for i := range comp {
		comp[i] = -1
	}
	compCnt := 0
	var dfs2 func(int)
	dfs2 = func(u int) {
		comp[u] = compCnt
		for _, p := range rg[u] {
			if comp[p] == -1 {
				dfs2(p)
			}
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == -1 {
			dfs2(v)
			compCnt++
		}
	}

	base := make([]int64, compCnt)
	dag := make([][]edge, compCnt)
	for _, e := range tc.edges {
		cu, cv := comp[e.u], comp[e.v]
		if cu == cv {
			base[cu] += gain(e.w)
		} else {
			dag[cu] = append(dag[cu], edge{u: cu, v: cv, w: e.w})
		}
	}

	memo := make([]int64, compCnt)
	seen := make([]bool, compCnt)
	var dp func(int) int64
	dp = func(c int) int64 {
		if seen[c] {
			return memo[c]
		}
		seen[c] = true
		best := int64(0)
		for _, e := range dag[c] {
			val := e.w + dp(e.v)
			if val > best {
				best = val
			}
		}
		memo[c] = base[c] + best
		return memo[c]
	}
	return dp(comp[tc.s])
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	m := rng.Intn(8)
	edges := make([]edge, 0, m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		w := int64(rng.Intn(30))
		edges = append(edges, edge{u: x, v: y, w: w})
	}
	s := rng.Intn(n) + 1
	return testCase{n: n, edges: edges, s: s}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	sb.WriteString(fmt.Sprintf("%d\n", tc.s))
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for t := 0; t < 100; t++ {
		tc := genCase(rng)
		input := buildInput(tc)
		exp := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: output is not an integer: %q\ninput:%s", t+1, got, input)
			os.Exit(1)
		}
		if gotVal != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:%s", t+1, exp, gotVal, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
