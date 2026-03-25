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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

// Embedded solver for 1984E
func solveTree(n int, edges [][2]int) int {
	if n == 1 {
		return 1
	}

	adj := make([][]int, n+1)
	for i := range adj {
		adj[i] = []int{}
	}
	degArr := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		degArr[u]++
		degArr[v]++
	}

	dp0 := make([]int, n+1) // max independent set not taking u
	dp1 := make([]int, n+1) // max independent set taking u

	// Iterative DFS for dp
	type frame struct {
		u, p  int
		childIdx int
	}
	stack := []frame{{1, 0, 0}}

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u, p := top.u, top.p

		if top.childIdx == 0 {
			dp0[u] = 0
			dp1[u] = 1
		}

		if top.childIdx < len(adj[u]) {
			v := adj[u][top.childIdx]
			top.childIdx++
			if v != p {
				stack = append(stack, frame{v, u, 0})
			}
			continue
		}

		// All children processed
		for _, v := range adj[u] {
			if v != p {
				if dp0[v] > dp1[v] {
					dp0[u] += dp0[v]
				} else {
					dp0[u] += dp1[v]
				}
				dp1[u] += dp0[v]
			}
		}

		stack = stack[:len(stack)-1]
	}

	alpha := dp0[1]
	if dp1[1] > dp0[1] {
		alpha = dp1[1]
	}

	// Reroot to check if any leaf has MIS == alpha when it's excluded
	// dfs2: check if there exists a leaf u such that the MIS of the tree
	// (with full rerooting) where u is not taken equals alpha
	found := false

	type frame2 struct {
		u, p     int
		out0, out1 int
		childIdx int
		s0, s1   int
	}

	// We need to compute for each node the "full" dp0/dp1 including parent contribution
	// S0[u] = sum over all neighbors (including parent contribution) of max(dp0[v], dp1[v])
	// S1[u] = 1 + sum over all neighbors of dp0[v]
	// For leaf u: check if S0[u] == alpha

	// Iterative rerooting
	type reroot struct {
		u, p, out0, out1 int
	}
	stack2 := []reroot{{1, 0, 0, 0}}
	for len(stack2) > 0 {
		cur := stack2[len(stack2)-1]
		stack2 = stack2[:len(stack2)-1]
		u, p, out0, out1 := cur.u, cur.p, cur.out0, cur.out1

		S0 := 0
		S1 := 1
		if p != 0 {
			if out0 > out1 {
				S0 += out0
			} else {
				S0 += out1
			}
			S1 += out0
		}
		for _, v := range adj[u] {
			if v != p {
				if dp0[v] > dp1[v] {
					S0 += dp0[v]
				} else {
					S0 += dp1[v]
				}
				S1 += dp0[v]
			}
		}

		if degArr[u] == 1 && S0 == alpha {
			found = true
			break
		}

		for _, v := range adj[u] {
			if v != p {
				vMax := dp0[v]
				if dp1[v] > dp0[v] {
					vMax = dp1[v]
				}
				nxtOut0 := S0 - vMax
				nxtOut1 := S1 - dp0[v]
				stack2 = append(stack2, reroot{v, u, nxtOut0, nxtOut1})
			}
		}
	}

	if found {
		return alpha + 1
	}
	return alpha
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 100; i++ {
		n := rng.Intn(8) + 2
		edges := randomTree(rng, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		input := fmt.Sprintf("1\n%s", sb.String())

		exp := fmt.Sprintf("%d", solveTree(n, edges))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
