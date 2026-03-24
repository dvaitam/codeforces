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

func generateCase(rng *rand.Rand) (int, int, [][2]int) {
	n := rng.Intn(10) + 1
	k := rng.Intn(n + 1)
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return n, k, edges
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// solve is an embedded reference solver for 1593E.
func solve(input string) string {
	r := strings.NewReader(input)
	var t int
	fmt.Fscan(r, &t)

	var outBuf strings.Builder

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(r, &n, &k)

		g := make([][]int, n+1)
		deg := make([]int, n+1)

		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(r, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
			deg[u]++
			deg[v]++
		}

		if k == 0 {
			fmt.Fprintln(&outBuf, n)
			continue
		}

		q := make([]int, 0)
		for i := 1; i <= n; i++ {
			if deg[i] <= 1 {
				q = append(q, i)
			}
		}

		rem := n
		for step := 0; step < k && len(q) > 0; step++ {
			rem -= len(q)
			nq := make([]int, 0)
			for _, v := range q {
				deg[v] = 0
				for _, to := range g[v] {
					if deg[to] > 0 {
						deg[to]--
						if deg[to] == 1 {
							nq = append(nq, to)
						}
					}
				}
			}
			q = nq
		}

		fmt.Fprintln(&outBuf, rem)
	}
	return strings.TrimSpace(outBuf.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, edges := generateCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		input := sb.String()
		exp := solve(input)
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
