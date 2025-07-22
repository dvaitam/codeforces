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

type op struct {
	t int
	v int
}

type testCase struct {
	n     int
	edges [][2]int
	ops   []op
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d\n", len(tc.ops))
	for _, op := range tc.ops {
		fmt.Fprintf(&sb, "%d %d\n", op.t, op.v)
	}
	input := sb.String()
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
	outputLines := strings.Fields(strings.TrimSpace(out.String()))
	expLines := expected(tc)
	if len(outputLines) != len(expLines) {
		return fmt.Errorf("expected %d lines got %d", len(expLines), len(outputLines))
	}
	for i := range expLines {
		if outputLines[i] != expLines[i] {
			return fmt.Errorf("line %d: expected %s got %s", i+1, expLines[i], outputLines[i])
		}
	}
	return nil
}

func expected(tc testCase) []string {
	n := tc.n
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n+1)
	order := []int{1}
	parent[1] = 0
	for i := 0; i < len(order); i++ {
		v := order[i]
		for _, to := range adj[v] {
			if to != parent[v] {
				parent[to] = v
				order = append(order, to)
			}
		}
	}
	filled := make([]bool, n+1)
	var res []string
	for _, op := range tc.ops {
		v := op.v
		switch op.t {
		case 1:
			// fill subtree of v
			stack := []int{v}
			for len(stack) > 0 {
				x := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if filled[x] {
					// already filled but continue to explore
				}
				filled[x] = true
				for _, to := range adj[x] {
					if to != parent[x] {
						stack = append(stack, to)
					}
				}
			}
		case 2:
			// empty path to root
			for x := v; x != 0; x = parent[x] {
				filled[x] = false
			}
		case 3:
			if filled[v] {
				res = append(res, "1")
			} else {
				res = append(res, "0")
			}
		}
	}
	return res
}

func generateCases(rng *rand.Rand) []testCase {
	cases := []testCase{}
	for len(cases) < 100 {
		n := rng.Intn(7) + 2
		edges := make([][2]int, n-1)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges[i-2] = [2]int{p, i}
		}
		q := rng.Intn(15) + 1
		ops := make([]op, q)
		for i := 0; i < q; i++ {
			t := rng.Intn(3) + 1
			v := rng.Intn(n) + 1
			ops[i] = op{t, v}
		}
		cases = append(cases, testCase{n, edges, ops})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateCases(rng)
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
