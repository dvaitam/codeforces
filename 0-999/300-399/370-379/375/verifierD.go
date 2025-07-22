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

type testCase struct{ input string }

func solveCase(in string) string {
	rdr := strings.NewReader(in)
	var n, m int
	fmt.Fscan(rdr, &n, &m)
	colors := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(rdr, &colors[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(rdr, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// build parent and subtree list
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v != parent[u] {
				parent[v] = u
				stack = append(stack, v)
			}
		}
	}
	// compute children
	children := make([][]int, n+1)
	for _, u := range order {
		if u == 1 {
			continue
		}
		children[parent[u]] = append(children[parent[u]], u)
	}
	// gather subtree nodes
	var collect func(int, []int) []int
	collect = func(u int, arr []int) []int {
		arr = append(arr, u)
		for _, v := range children[u] {
			arr = collect(v, arr)
		}
		return arr
	}
	var sb strings.Builder
	for i := 0; i < m; i++ {
		var v, k int
		fmt.Fscan(rdr, &v, &k)
		nodes := collect(v, nil)
		freq := map[int]int{}
		for _, x := range nodes {
			freq[colors[x]]++
		}
		ans := 0
		for _, c := range freq {
			if c >= k {
				ans++
			}
		}
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprint(ans))
	}
	return sb.String()
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
	got := strings.TrimSpace(out.String())
	exp := solveCase(tc.input)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 1; i <= n; i++ {
		color := rng.Intn(5) + 1
		sb.WriteString(fmt.Sprintf("%d ", color))
	}
	sb.WriteByte('\n')
	edges := randomTree(rng, n)
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 0; i < m; i++ {
		v := rng.Intn(n) + 1
		k := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", v, k))
	}
	return testCase{input: sb.String()}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 105)
	cases = append(cases, randomCase(rng))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
