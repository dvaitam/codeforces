package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(n, m, s int, edges [][2]int) int {
	adj := make([][]int, n+1)
	rev := make([][]int, n+1)

	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		rev[v] = append(rev[v], u)
	}

	visited := make([]bool, n+1)
	order := make([]int, 0, n)

	var dfs1 func(int)
	dfs1 = func(u int) {
		visited[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				dfs1(v)
			}
		}
		order = append(order, u)
	}

	for i := 1; i <= n; i++ {
		if !visited[i] {
			dfs1(i)
		}
	}

	scc := make([]int, n+1)
	for i := 1; i <= n; i++ {
		visited[i] = false
	}

	var dfs2 func(int, int)
	dfs2 = func(u, c int) {
		visited[u] = true
		scc[u] = c
		for _, v := range rev[u] {
			if !visited[v] {
				dfs2(v, c)
			}
		}
	}

	count := 0
	for i := n - 1; i >= 0; i-- {
		u := order[i]
		if !visited[u] {
			count++
			dfs2(u, count)
		}
	}

	inDegree := make([]int, count+1)
	for u := 1; u <= n; u++ {
		for _, v := range adj[u] {
			if scc[u] != scc[v] {
				inDegree[scc[v]]++
			}
		}
	}

	ans := 0
	for i := 1; i <= count; i++ {
		if inDegree[i] == 0 && i != scc[s] {
			ans++
		}
	}

	return ans
}

func runBinary(bin, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), stderr.String(), err
}

type testCase struct{ input string }

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(n*n + 1)
	s := rng.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, s))
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
	}
	return testCase{sb.String()}
}

func parseCase(tc testCase) (int, int, int, [][2]int) {
	r := strings.NewReader(tc.input)
	var n, m, s int
	fmt.Fscan(r, &n, &m, &s)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &edges[i][0], &edges[i][1])
	}
	return n, m, s, edges
}

func runCase(bin string, tc testCase, idx int) error {
	n, _, s, edges := parseCase(tc)
	expected := solve(n, 0, s, edges)
	gotOut, gotErr, err := runBinary(bin, tc.input)
	if err != nil {
		return fmt.Errorf("test %d: runtime error: %v\nstderr: %s", idx, err, gotErr)
	}
	if strings.TrimSpace(gotOut) != fmt.Sprint(expected) {
		return fmt.Errorf("test %d failed\nexpected: %d\n got: %s", idx, expected, gotOut)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(99))
	for i := 1; i <= 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc, i); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
