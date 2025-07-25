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

func solveCase(n int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stk := []int{1}
	parent[1] = 0
	for len(stk) > 0 {
		v := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		order = append(order, v)
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			stk = append(stk, u)
		}
	}
	sz := make([]int, n+1)
	childMax1 := make([]int, n+1)
	childMax2 := make([]int, n+1)
	childMaxNode := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		sz[v] = 1
		m1, m2, mn := 0, 0, 0
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			sz[v] += sz[u]
			if sz[u] > m1 {
				m2 = m1
				m1 = sz[u]
				mn = u
			} else if sz[u] > m2 {
				m2 = sz[u]
			}
		}
		childMax1[v] = m1
		childMax2[v] = m2
		childMaxNode[v] = mn
	}
	upSz := make([]int, n+1)
	compMax1 := make([]int, n+1)
	compMax2 := make([]int, n+1)
	compMaxNode := make([]int, n+1)
	half := n / 2
	for _, v := range order {
		if v == 1 {
			upSz[v] = 0
		} else {
			upSz[v] = n - sz[v]
		}
		m1, m2, mn := childMax1[v], childMax2[v], childMaxNode[v]
		if upSz[v] > m1 {
			m2 = m1
			m1 = upSz[v]
			mn = parent[v]
		} else if upSz[v] > m2 {
			m2 = upSz[v]
		}
		compMax1[v] = m1
		compMax2[v] = m2
		compMaxNode[v] = mn
	}
	out := make([]byte, n)
	for v := 1; v <= n; v++ {
		m := compMax1[v]
		if m <= half {
			out[v-1] = '1'
			continue
		}
		w := compMaxNode[v]
		var local int
		if compMaxNode[w] != v {
			local = compMax1[w]
		} else {
			local = compMax2[w]
		}
		if 2*local >= m {
			out[v-1] = '1'
		} else {
			out[v-1] = '0'
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte(out[i])
	}
	return sb.String()
}

func runCase(bin string, n int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := solveCase(n, edges)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func randomTree(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(20) + 2
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return n, edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, edges := randomTree(rng)
		if err := runCase(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
