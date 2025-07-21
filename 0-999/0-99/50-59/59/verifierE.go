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

func solveCaseE(n int, edges [][2]int, forb [][3]int) (string, string) {
	adj := make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	n1 := int64(n) + 1
	forbid := make(map[int64]bool)
	for _, t := range forb {
		a, b, c := t[0], t[1], t[2]
		key := (int64(a)*n1+int64(b))*n1 + int64(c)
		forbid[key] = true
	}
	type state struct{ prev, cur int }
	start := state{0, 1}
	key := func(s state) int64 { return int64(s.prev)*n1 + int64(s.cur) }
	dist := map[int64]int{key(start): 0}
	parent := make(map[int64]int64)
	q := []state{start}
	var end state
	found := false
	for len(q) > 0 && !found {
		s := q[0]
		q = q[1:]
		if s.cur == n {
			end = s
			found = true
			break
		}
		for _, w := range adj[s.cur] {
			if s.prev != 0 {
				if forbid[(int64(s.prev)*n1+int64(s.cur))*n1+int64(w)] {
					continue
				}
			}
			ns := state{s.cur, w}
			nk := key(ns)
			if _, ok := dist[nk]; !ok {
				dist[nk] = dist[key(s)] + 1
				parent[nk] = key(s)
				q = append(q, ns)
			}
		}
	}
	if !found {
		return "-1", ""
	}
	path := []int{end.cur}
	for k := parent[key(end)]; k != -1; k = parent[k] {
		prev := int(k % n1)
		path = append(path, prev)
		if prev == 1 && int(k/n1) == 0 {
			break
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	d := len(path) - 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", d))
	for i, v := range path {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String(), ""
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	nodes := make([]int, n)
	for i := range nodes {
		nodes[i] = i + 1
	}
	edges := make([][2]int, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if rng.Intn(2) == 0 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	if len(edges) == 0 {
		edges = append(edges, [2]int{1, n})
	}
	forb := make([][3]int, 0)
	maxForb := rng.Intn(3)
	for i := 0; i < maxForb; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		c := rng.Intn(n) + 1
		for a == b {
			b = rng.Intn(n) + 1
		}
		for c == a || c == b {
			c = rng.Intn(n) + 1
		}
		forb = append(forb, [3]int{a, b, c})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), len(forb)))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for _, t := range forb {
		fmt.Fprintf(&sb, "%d %d %d\n", t[0], t[1], t[2])
	}
	expect, _ := solveCaseE(n, edges, forb)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
