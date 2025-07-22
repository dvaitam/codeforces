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

var (
	k       int
	adj     [][]int
	removed []bool
	sz      []int
	ans     int64
)

func dfsSize(u, p int) {
	sz[u] = 1
	for _, v := range adj[u] {
		if v != p && !removed[v] {
			dfsSize(v, u)
			sz[u] += sz[v]
		}
	}
}

func dfsCentroid(u, p, total int) int {
	for _, v := range adj[u] {
		if v != p && !removed[v] {
			if sz[v] > total/2 {
				return dfsCentroid(v, u, total)
			}
		}
	}
	return u
}

func dfsCount(u, p, depth int, cnt []int) {
	if depth > k {
		return
	}
	cnt[depth]++
	for _, v := range adj[u] {
		if v != p && !removed[v] {
			dfsCount(v, u, depth+1, cnt)
		}
	}
}

func solve(u int) {
	dfsSize(u, -1)
	c := dfsCentroid(u, -1, sz[u])
	removed[c] = true
	freq := make([]int, k+1)
	freq[0] = 1
	for _, v := range adj[c] {
		if removed[v] {
			continue
		}
		cnt := make([]int, k+1)
		dfsCount(v, c, 1, cnt)
		for d := 1; d <= k; d++ {
			if cnt[d] > 0 && k-d >= 0 {
				ans += int64(cnt[d]) * int64(freq[k-d])
			}
		}
		for d := 1; d <= k; d++ {
			freq[d] += cnt[d]
		}
	}
	for _, v := range adj[c] {
		if !removed[v] {
			solve(v)
		}
	}
}

type Test struct {
	n     int
	k     int
	edges [][2]int
}

func generateTest() Test {
	n := rand.Intn(10) + 1
	k := rand.Intn(n)
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return Test{n, k, edges}
}

func (t Test) Input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", t.n, t.k)
	for _, e := range t.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func solveTest(t Test) int64 {
	k = t.k
	adj = make([][]int, t.n+1)
	removed = make([]bool, t.n+1)
	sz = make([]int, t.n+1)
	for _, e := range t.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	ans = 0
	solve(1)
	return ans
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		t := generateTest()
		inp := t.Input()
		exp := solveTest(t)
		out, err := runBinary(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int64
		if _, e := fmt.Sscan(out, &got); e != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\n", i+1, e)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:\n%s\n", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
