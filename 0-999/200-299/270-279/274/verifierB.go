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

type edge struct{ u, v int }

func expectedB(n int, edges []edge, vals []int64) int64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	var dfs func(int, int) (int64, int64)
	dfs = func(u, p int) (int64, int64) {
		var inc, dec int64
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			ci, cd := dfs(v, u)
			if ci > inc {
				inc = ci
			}
			if cd > dec {
				dec = cd
			}
		}
		cur := vals[u] + inc - dec
		if cur > 0 {
			dec += cur
		} else {
			inc += -cur
		}
		return inc, dec
	}
	inc, dec := dfs(1, 0)
	return inc + dec
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	vals := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		vals[i] = int64(rng.Intn(11) - 5) // -5..5
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n) + "\n")
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(vals[i], 10))
	}
	sb.WriteByte('\n')
	exp := expectedB(n, edges, vals)
	return sb.String(), fmt.Sprintf("%d\n", exp)
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
