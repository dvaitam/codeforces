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

type Edge struct{ u, v int }

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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(n int, a []int, edges []Edge, queries [][2]int) []int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	res := make([]int, len(queries))
	for qi, q := range queries {
		u, v := q[0], q[1]
		ans := 0
		dist := 0
		cur := v
		for cur != u {
			val := a[cur] ^ dist
			if val > ans {
				ans = val
			}
			cur = parent[cur]
			dist++
		}
		val := a[cur] ^ dist
		if val > ans {
			ans = val
		}
		res[qi] = ans
	}
	return res
}

func genCase(rng *rand.Rand) (string, int, []int, []Edge, [][2]int) {
	n := rng.Intn(8) + 2
	q := rng.Intn(10) + 1
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(20)
	}
	edges := make([]Edge, 0, n-1)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		parent[i] = p
		edges = append(edges, Edge{p, i})
	}
	queries := make([][2]int, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	for i := 0; i < q; i++ {
		v := rng.Intn(n) + 1
		// pick ancestor
		cur := v
		path := []int{cur}
		for cur != 1 {
			cur = parent[cur]
			path = append(path, cur)
		}
		u := path[rng.Intn(len(path))]
		queries[i] = [2]int{u, v}
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
	}
	return sb.String(), n, a[1:], edges, queries
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, a, edges, qs := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != len(qs) {
			fmt.Fprintf(os.Stderr, "case %d expected %d lines got %d\n", i+1, len(qs), len(lines))
			os.Exit(1)
		}
		exp := expected(n, a, edges, qs)
		for j, line := range lines {
			v, err := strconv.Atoi(strings.TrimSpace(line))
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d invalid int on line %d\n", i+1, j+1)
				os.Exit(1)
			}
			if v != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d query %d failed: expected %d got %d\ninput:\n%s", i+1, j+1, exp[j], v, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
