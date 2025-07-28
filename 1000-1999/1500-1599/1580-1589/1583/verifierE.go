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

const LOGE = 20

type caseE struct {
	n, m    int
	edges   [][2]int
	q       int
	queries [][2]int
}

func genCase(rng *rand.Rand) caseE {
	n := rng.Intn(6) + 2
	m := n - 1
	edges := make([][2]int, 0, m)
	// create random tree
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		for a == b {
			b = rng.Intn(n) + 1
		}
		queries[i] = [2]int{a, b}
	}
	return caseE{n, m, edges, q, queries}
}

func lca(u, v int, depth []int, up [][]int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for k := LOGE - 1; k >= 0; k-- {
		if diff&(1<<k) != 0 {
			u = up[k][u]
		}
	}
	if u == v {
		return u
	}
	for k := LOGE - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			u = up[k][u]
			v = up[k][v]
		}
	}
	return up[0][u]
}

func getPath(a, b int, parent []int, depth []int, up [][]int) []int {
	l := lca(a, b, depth, up)
	var path []int
	u := a
	for u != l {
		path = append(path, u)
		u = parent[u]
	}
	path = append(path, l)
	var tail []int
	v := b
	for v != l {
		tail = append(tail, v)
		v = parent[v]
	}
	for i := len(tail) - 1; i >= 0; i-- {
		path = append(path, tail[i])
	}
	return path
}

func expected(tc caseE) (string, []string) {
	g := make([][]int, tc.n+1)
	for _, e := range tc.edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	parent := make([]int, tc.n+1)
	depth := make([]int, tc.n+1)
	vis := make([]bool, tc.n+1)
	q := []int{1}
	vis[1] = true
	for i := 0; i < len(q); i++ {
		v := q[i]
		for _, to := range g[v] {
			if !vis[to] {
				vis[to] = true
				parent[to] = v
				depth[to] = depth[v] + 1
				q = append(q, to)
			}
		}
	}
	up := make([][]int, LOGE)
	for i := range up {
		up[i] = make([]int, tc.n+1)
	}
	for v := 1; v <= tc.n; v++ {
		up[0][v] = parent[v]
	}
	for k := 1; k < LOGE; k++ {
		for v := 1; v <= tc.n; v++ {
			up[k][v] = up[k-1][up[k-1][v]]
		}
	}

	parity := make([]int, tc.n+1)
	for _, qu := range tc.queries {
		parity[qu[0]] ^= 1
		parity[qu[1]] ^= 1
	}
	odd := 0
	for i := 1; i <= tc.n; i++ {
		if parity[i] == 1 {
			odd++
		}
	}
	if odd > 0 {
		return "NO", []string{fmt.Sprintf("%d", odd/2)}
	}
	out := []string{"YES"}
	for _, qu := range tc.queries {
		p := getPath(qu[0], qu[1], parent, depth, up)
		out = append(out, fmt.Sprintf("%d", len(p)))
		var sb strings.Builder
		for i, node := range p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", node))
		}
		out = append(out, sb.String())
	}
	return "", out
}

func runCase(bin string, tc caseE) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", tc.q))
	for _, qu := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outFields := strings.Fields(out.String())
	expPrefix, exp := expected(tc)
	if expPrefix == "NO" {
		if len(outFields) < 2 {
			return fmt.Errorf("expected NO and number")
		}
		if strings.ToUpper(outFields[0]) != "NO" {
			return fmt.Errorf("expected NO got %s", outFields[0])
		}
		if outFields[1] != exp[0] {
			return fmt.Errorf("expected %s got %s", exp[0], outFields[1])
		}
		return nil
	}
	// YES case
	if len(outFields) < 1 {
		return fmt.Errorf("no output")
	}
	if strings.ToUpper(outFields[0]) != "YES" {
		return fmt.Errorf("expected YES got %s", outFields[0])
	}
	outFields = outFields[1:]
	if len(outFields) != len(exp)-1 {
		return fmt.Errorf("expected %d tokens got %d", len(exp)-1, len(outFields))
	}
	for i, token := range outFields {
		if token != exp[i+1] {
			return fmt.Errorf("token %d expected %s got %s", i+1, exp[i+1], token)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
