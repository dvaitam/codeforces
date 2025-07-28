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

type edge struct{ u, v int }

type testCaseE struct {
	n     int
	vals  []int
	edges []edge
}

func checkRoot(adj [][]int, val []int, root int) bool {
	freq := make(map[int]int)
	stack := []struct {
		node, parent int
		enter        bool
	}{{root, -1, true}}
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if s.enter {
			if freq[val[s.node]] >= 1 {
				return false
			}
			freq[val[s.node]]++
			stack = append(stack, struct {
				node, parent int
				enter        bool
			}{s.node, s.parent, false})
			for _, to := range adj[s.node] {
				if to == s.parent {
					continue
				}
				stack = append(stack, struct {
					node, parent int
					enter        bool
				}{to, s.node, true})
			}
		} else {
			freq[val[s.node]]--
		}
	}
	return true
}

func solveE(tc testCaseE) int {
	adj := make([][]int, tc.n+1)
	for _, e := range tc.edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	count := 0
	val := append([]int{0}, tc.vals...)
	for r := 1; r <= tc.n; r++ {
		if checkRoot(adj, val, r) {
			count++
		}
	}
	return count
}

func buildInputE(tc testCaseE) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String()
}

func runCaseE(bin string, tc testCaseE) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputE(tc))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := solveE(tc)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCasesE() []testCaseE {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseE, 0, 100)
	for len(cases) < 100 {
		n := rng.Intn(7) + 1
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			vals[i] = rng.Intn(10) + 1
		}
		edges := make([]edge, n-1)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges[i-2] = edge{p, i}
		}
		cases = append(cases, testCaseE{n, vals, edges})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesE()
	for i, tc := range cases {
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
