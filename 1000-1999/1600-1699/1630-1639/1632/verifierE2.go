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

func solveE2(n int, edges [][2]int) []int {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	d := make([]int, n)
	maxH := 0
	var dfs func(int, int, int) int
	dfs = func(u, p, dep int) int {
		h1, h2 := dep, dep
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			q := dfs(v, u, dep+1)
			if q > h1 {
				h2 = h1
				h1 = q
			} else if q > h2 {
				h2 = q
			}
		}
		if h1 > maxH {
			maxH = h1
		}
		x := h1
		if h2 < x {
			x = h2
		}
		x--
		if x >= 0 {
			val := h1 + h2 - 2*dep + 1
			if val > d[x] {
				d[x] = val
			}
		}
		return h1
	}
	dfs(0, -1, 0)
	for i := n - 2; i >= 0; i-- {
		if d[i+1] > d[i] {
			d[i] = d[i+1]
		}
	}
	ans := 0
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		for ans < maxH && d[ans]/2+i > ans {
			ans++
		}
		res[i-1] = ans
	}
	return res
}

func randTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{i, p})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(8) + 2
	edges := randTree(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	ans := solveE2(n, edges)
	return sb.String(), ans
}

func runCase(bin, input string, exp []int) error {
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
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(parts))
	}
	for i, p := range parts {
		var v int
		if _, err := fmt.Sscan(p, &v); err != nil {
			return fmt.Errorf("bad int at pos %d: %v", i+1, err)
		}
		if v != exp[i] {
			return fmt.Errorf("pos %d expected %d got %d", i+1, exp[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
