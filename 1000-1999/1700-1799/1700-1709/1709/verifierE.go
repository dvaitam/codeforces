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

type pair struct{ a, b int }

func solveCase(n int, vals []int, edges []pair) int {
	g := make([][]int, n)
	for _, e := range edges {
		u, v := e.a, e.b
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	px := make([]int, n)
	var dfsPX func(v, p int)
	dfsPX = func(v, p int) {
		for _, to := range g[v] {
			if to == p {
				continue
			}
			px[to] = px[v] ^ vals[to]
			dfsPX(to, v)
		}
	}
	px[0] = vals[0]
	dfsPX(0, -1)
	ans := 0
	var dfs func(v, p int) map[int]struct{}
	dfs = func(v, p int) map[int]struct{} {
		set := map[int]struct{}{px[v]: {}}
		for _, to := range g[v] {
			if to == p {
				continue
			}
			child := dfs(to, v)
			if len(set) < len(child) {
				set, child = child, set
			}
			intersect := false
			for val := range child {
				if _, ok := set[val]; ok {
					intersect = true
					break
				}
			}
			if intersect {
				ans++
			} else {
				for val := range child {
					set[val] = struct{}{}
				}
			}
		}
		return set
	}
	dfs(0, -1)
	return ans
}

func genTree(rng *rand.Rand, n int) []pair {
	edges := make([]pair, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, pair{p, i})
	}
	return edges
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(7) + 1 // 1..7
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(16) + 1
	}
	edges := genTree(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", vals[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.a+1, e.b+1))
	}
	expect := solveCase(n, vals, edges)
	return sb.String(), expect
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != fmt.Sprintf("%d", expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, expect, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
