package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type edge struct{ u, v int }

type testCase struct {
	n     int
	edges []edge
}

func solve(tc testCase) string {
	n := tc.n
	adj := make([]int, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = -1
	}
	apr := make([]int, n+1)
	tem := []int{}
	to := []int{}
	nxt := []int{}
	vis := []bool{}
	vise := []bool{}
	cntege := 0
	var addDir func(u, v int)
	addDir = func(u, v int) {
		e := len(to)
		to = append(to, v)
		nxt = append(nxt, adj[u])
		adj[u] = e
		vis = append(vis, false)
		vise = append(vise, false)
	}
	addEdge := func(u, v int) {
		addDir(u, v)
		addDir(v, u)
		cntege++
	}
	for _, e := range tc.edges {
		addEdge(e.u, e.v)
		apr[e.u]++
		apr[e.v]++
	}
	for i := 1; i <= n; i++ {
		if apr[i]%2 != 0 {
			tem = append(tem, i)
		}
	}
	for i := 0; i < len(tem); i += 2 {
		if i+1 < len(tem) {
			addEdge(tem[i], tem[i+1])
		} else {
			addEdge(tem[i], tem[i])
		}
	}
	iter := make([]int, n+1)
	copy(iter, adj)
	nodeSt := []int{1}
	edgeSt := []int{-1}
	res := []int{}
	for len(nodeSt) > 0 {
		v := nodeSt[len(nodeSt)-1]
		if iter[v] != -1 {
			e := iter[v]
			iter[v] = nxt[e]
			if vis[e] {
				continue
			}
			vis[e] = true
			vis[e^1] = true
			nodeSt = append(nodeSt, to[e])
			edgeSt = append(edgeSt, e)
		} else {
			nodeSt = nodeSt[:len(nodeSt)-1]
			e := edgeSt[len(edgeSt)-1]
			edgeSt = edgeSt[:len(edgeSt)-1]
			if e != -1 {
				res = append(res, e)
			}
		}
	}
	z := 0
	for _, e := range res {
		z++
		if z%2 == 1 {
			vise[e] = true
		} else {
			vise[e^1] = true
		}
	}
	if cntege%2 != 0 {
		addDir(1, 1)
		e2 := len(to)
		addDir(1, 1)
		cntege++
		vise[e2] = true
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, cntege)
	for u := 1; u <= n; u++ {
		for e := adj[u]; e != -1; e = nxt[e] {
			if vise[e] {
				fmt.Fprintf(&sb, "%d %d\n", u, to[e])
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(6) + 2
		m := rnd.Intn(8) + 1
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			u := rnd.Intn(n) + 1
			v := rnd.Intn(n) + 1
			edges[i] = edge{u, v}
		}
		tests = append(tests, testCase{n, edges})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
		input := sb.String()
		expected := solve(tc)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		want := strings.TrimSpace(expected)
		if got != want {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\n got:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
