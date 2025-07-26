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

func solveE(n int, h []int64, edges [][2]int) (string, string) {
	gRev := make([][]int, n)
	outdeg := make([]int, n)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		gRev[v] = append(gRev[v], u)
		outdeg[u]++
	}
	dp := make([]int, n)
	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if outdeg[i] == 0 {
			q = append(q, i)
		}
	}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, p := range gRev[u] {
			if dp[p] < dp[u]+1 {
				dp[p] = dp[u] + 1
			}
			outdeg[p]--
			if outdeg[p] == 0 {
				q = append(q, p)
			}
		}
	}
	var x int64
	even := make([]bool, n)
	for i := 0; i < n; i++ {
		if dp[i]%2 == 0 {
			x ^= h[i]
			even[i] = true
		}
	}
	if x == 0 {
		var sb strings.Builder
		sb.WriteString("LOSE\n")
		return sb.String(), ""
	}
	newH := make([]int64, n)
	copy(newH, h)
	for i := 0; i < n; i++ {
		if even[i] {
			want := h[i] ^ x
			if want < h[i] {
				newH[i] = want
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("WIN\n")
	for i, v := range newH {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), ""
}

func genCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(n*(n-1)/2 + 1)
	edges := make([][2]int, 0, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		for u == v || hasEdge(edges, u, v) || createsCycle(edges, u, v) {
			u = rng.Intn(n)
			v = rng.Intn(n)
		}
		edges = append(edges, [2]int{u, v})
	}
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		h[i] = int64(rng.Intn(10))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	expect, _ := solveE(n, h, edges)
	return sb.String(), expect
}

func hasEdge(edges [][2]int, u, v int) bool {
	for _, e := range edges {
		if e[0] == u && e[1] == v {
			return true
		}
	}
	return false
}
func createsCycle(edges [][2]int, u, v int) bool { // simple check via DFS
	n := 0
	for _, e := range edges {
		if e[0] >= n {
			n = e[0] + 1
		}
		if e[1] >= n {
			n = e[1] + 1
		}
	}
	n = max(n, max(u+1, v+1))
	g := make([][]int, n)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
	}
	vis := make([]int, n)
	var dfs func(int) bool
	dfs = func(x int) bool {
		vis[x] = 1
		for _, to := range g[x] {
			if vis[to] == 1 {
				return true
			}
			if vis[to] == 0 {
				if dfs(to) {
					return true
				}
			}
		}
		vis[x] = 2
		return false
	}
	g[u] = append(g[u], v)
	for i := 0; i < n; i++ {
		vis[i] = 0
	}
	for i := 0; i < n; i++ {
		if vis[i] == 0 {
			if dfs(i) {
				return true
			}
		}
	}
	return false
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in, expect := genCaseE(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
