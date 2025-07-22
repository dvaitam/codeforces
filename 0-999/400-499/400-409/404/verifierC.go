package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseC struct {
	n int
	k int
	d []int
}

type edge struct{ u, v int }

func generateValidC(r *rand.Rand) testCaseC {
	for {
		n := 3 + r.Intn(6)
		k := 1 + r.Intn(3)
		parents := make([]int, n+1)
		depth := make([]int, n+1)
		deg := make([]int, n+1)
		ok := true
		for i := 2; i <= n; i++ {
			cand := []int{}
			for p := 1; p < i; p++ {
				if deg[p] < k {
					cand = append(cand, p)
				}
			}
			if len(cand) == 0 {
				ok = false
				break
			}
			par := cand[r.Intn(len(cand))]
			parents[i] = par
			depth[i] = depth[par] + 1
			deg[par]++
		}
		if !ok {
			continue
		}
		d := make([]int, n)
		for i := 1; i <= n; i++ {
			d[i-1] = depth[i]
		}
		return testCaseC{n: n, k: k, d: d}
	}
}

func generateInvalidC(r *rand.Rand) testCaseC {
	n := 3 + r.Intn(6)
	k := 1 + r.Intn(3)
	d := make([]int, n)
	for i := 0; i < n; i++ {
		d[i] = r.Intn(n)
	}
	if r.Intn(2) == 0 && n > 1 {
		d[0] = 0
		d[1] = 0
	}
	return testCaseC{n: n, k: k, d: d}
}

func generateTestsC() []testCaseC {
	r := rand.New(rand.NewSource(3))
	tests := make([]testCaseC, 0, 100)
	for len(tests) < 50 {
		tests = append(tests, generateValidC(r))
	}
	for len(tests) < 100 {
		tests = append(tests, generateInvalidC(r))
	}
	return tests
}

func solveC(t testCaseC) (bool, []edge) {
	n, k := t.n, t.k
	d := make([]int, n+1)
	levels := make([][]int, n)
	maxD := 0
	for i := 1; i <= n; i++ {
		di := t.d[i-1]
		if di < 0 || di >= n {
			return false, nil
		}
		d[i] = di
		levels[di] = append(levels[di], i)
		if di > maxD {
			maxD = di
		}
	}
	if len(levels[0]) != 1 {
		return false, nil
	}
	for dist := 1; dist <= maxD; dist++ {
		cnt := len(levels[dist])
		if cnt == 0 {
			continue
		}
		parentCnt := len(levels[dist-1])
		cap := parentCnt * (k - 1)
		if dist == 1 {
			cap = parentCnt * k
		}
		if cnt > cap {
			return false, nil
		}
	}
	edges := make([]edge, 0, n-1)
	if maxD >= 1 {
		root := levels[0][0]
		for _, v := range levels[1] {
			edges = append(edges, edge{root, v})
		}
	}
	for dist := 2; dist <= maxD; dist++ {
		parents := levels[dist-1]
		children := levels[dist]
		if len(children) == 0 {
			continue
		}
		count := make([]int, len(parents))
		idx := 0
		for _, v := range children {
			for idx < len(parents) && count[idx] >= k-1 {
				idx++
			}
			if idx >= len(parents) {
				return false, nil
			}
			u := parents[idx]
			edges = append(edges, edge{u, v})
			count[idx]++
		}
	}
	return true, edges
}

func validateOutput(t testCaseC, out string) bool {
	trimmed := strings.TrimSpace(out)
	if trimmed == "-1" {
		ok, _ := solveC(t)
		return !ok
	}
	in := bufio.NewReader(strings.NewReader(out))
	var m int
	if _, err := fmt.Fscan(in, &m); err != nil {
		return false
	}
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(in, &edges[i].u, &edges[i].v); err != nil {
			return false
		}
	}
	_, _ = io.ReadAll(in)
	return checkGraph(t, edges)
}

func checkGraph(t testCaseC, edges []edge) bool {
	n, k := t.n, t.k
	root := -1
	for i, di := range t.d {
		if di == 0 {
			if root != -1 {
				return false
			}
			root = i + 1
		}
	}
	if root == -1 {
		return false
	}
	deg := make([]int, n+1)
	adj := make([][]int, n+1)
	seen := make(map[[2]int]bool)
	for _, e := range edges {
		if e.u < 1 || e.u > n || e.v < 1 || e.v > n || e.u == e.v {
			return false
		}
		p := [2]int{e.u, e.v}
		if p[0] > p[1] {
			p[0], p[1] = p[1], p[0]
		}
		if seen[p] {
			return false
		}
		seen[p] = true
		deg[e.u]++
		deg[e.v]++
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	for i := 1; i <= n; i++ {
		if deg[i] > k {
			return false
		}
	}
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{root}
	dist[root] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if dist[i] != t.d[i-1] {
			return false
		}
	}
	return true
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsC()
	for i, t := range tests {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
		for idx, di := range t.d {
			if idx > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(fmt.Sprint(di))
		}
		b.WriteByte('\n')
		out, err := runBinary(bin, b.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if !validateOutput(t, out) {
			fmt.Printf("test %d failed\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
