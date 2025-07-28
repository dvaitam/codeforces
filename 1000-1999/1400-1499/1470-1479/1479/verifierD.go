package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type queryD struct{ u, v, l, r int }

type testCaseD struct {
	n       int
	q       int
	a       []int
	edges   [][2]int
	queries []queryD
}

func genTestsD() []testCaseD {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rng.Intn(6) + 2 // 2..7
		a := make([]int, n+1)
		for j := 1; j <= n; j++ {
			a[j] = rng.Intn(n) + 1
		}
		edges := make([][2]int, 0, n-1)
		for j := 2; j <= n; j++ {
			p := rng.Intn(j-1) + 1
			edges = append(edges, [2]int{p, j})
		}
		q := rng.Intn(5) + 1
		qs := make([]queryD, q)
		for j := 0; j < q; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			l := rng.Intn(n) + 1
			r := l + rng.Intn(n+1-l)
			qs[j] = queryD{u, v, l, r}
		}
		tests[i] = testCaseD{n: n, q: q, a: a, edges: edges, queries: qs}
	}
	return tests
}

func findPath(n int, edges [][2]int, u, v int) []int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	parent := make([]int, n+1)
	vis := make([]bool, n+1)
	q := []int{u}
	vis[u] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur == v {
			break
		}
		for _, nxt := range adj[cur] {
			if !vis[nxt] {
				vis[nxt] = true
				parent[nxt] = cur
				q = append(q, nxt)
			}
		}
	}
	path := []int{v}
	for v != u {
		v = parent[v]
		path = append(path, v)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func expectedAns(tc testCaseD) []int {
	results := make([]int, tc.q)
	for i, qu := range tc.queries {
		p := findPath(tc.n, tc.edges, qu.u, qu.v)
		freq := make(map[int]int)
		for _, node := range p {
			val := tc.a[node]
			freq[val]++
		}
		ans := -1
		for x := qu.l; x <= qu.r; x++ {
			if freq[x]%2 == 1 {
				ans = x
				break
			}
		}
		results[i] = ans
	}
	return results
}

func buildInput(tc testCaseD) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, tc.a[i])
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for _, q := range tc.queries {
		fmt.Fprintf(&sb, "%d %d %d %d\n", q.u, q.v, q.l, q.r)
	}
	return sb.String()
}

func runCase(bin string, tc testCaseD, expect []int) error {
	input := buildInput(tc)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(expect) {
		return fmt.Errorf("expected %d outputs got %d", len(expect), len(fields))
	}
	for i, exp := range expect {
		val, err := strconv.Atoi(fields[i])
		if err != nil || val != exp {
			return fmt.Errorf("on query %d expected %d got %s", i+1, exp, fields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTestsD()
	for i, tc := range cases {
		exp := expectedAns(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
