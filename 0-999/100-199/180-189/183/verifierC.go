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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func computeExpected(n int, edges [][2]int) int {
	g := make([][]int, n)
	grev := make([][]int, n)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		g[u] = append(g[u], v)
		grev[v] = append(grev[v], u)
	}
	visited := make([]bool, n)
	order := make([]int, 0, n)
	type frame struct{ v, i int }
	for s := 0; s < n; s++ {
		if visited[s] {
			continue
		}
		stack := []frame{{s, 0}}
		for len(stack) > 0 {
			fr := &stack[len(stack)-1]
			v, i := fr.v, fr.i
			if i == 0 {
				visited[v] = true
			}
			if fr.i < len(g[v]) {
				to := g[v][fr.i]
				fr.i++
				if !visited[to] {
					stack = append(stack, frame{to, 0})
				}
			} else {
				order = append(order, v)
				stack = stack[:len(stack)-1]
			}
		}
	}
	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	cid := 0
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] != -1 {
			continue
		}
		stack := []int{v}
		comp[v] = cid
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, w := range grev[u] {
				if comp[w] == -1 {
					comp[w] = cid
					stack = append(stack, w)
				}
			}
		}
		cid++
	}
	nodes := make([][]int, cid)
	for v := 0; v < n; v++ {
		nodes[comp[v]] = append(nodes[comp[v]], v)
	}
	globalG := 0
	depth := make([]int, n)
	mark := make([]bool, n)
	for c := 0; c < cid; c++ {
		if len(nodes[c]) == 0 {
			continue
		}
		for _, v := range nodes[c] {
			mark[v] = false
			depth[v] = 0
		}
		root := nodes[c][0]
		mark[root] = true
		depth[root] = 0
		stack := []int{root}
		localG := 0
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, v := range g[u] {
				if comp[v] != c {
					continue
				}
				if !mark[v] {
					mark[v] = true
					depth[v] = depth[u] + 1
					stack = append(stack, v)
				} else {
					d := depth[u] + 1 - depth[v]
					if d < 0 {
						d = -d
					}
					localG = gcd(localG, d)
				}
			}
		}
		globalG = gcd(globalG, localG)
	}
	if globalG == 0 {
		globalG = n
	}
	return globalG
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	m := rng.Intn(15)
	edges := make([][2]int, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		edges[i] = [2]int{u, v}
		fmt.Fprintf(&sb, "%d %d\n", u+1, v+1)
	}
	ans := computeExpected(n, edges)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
