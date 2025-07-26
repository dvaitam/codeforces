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

type Edge struct{ u, v int }

func isBipartiteSub(n int, edges []Edge, L, R int) bool {
	color := make(map[int]int)
	for v := L; v <= R; v++ {
		if _, ok := color[v]; !ok {
			queue := []int{v}
			color[v] = 0
			for len(queue) > 0 {
				x := queue[0]
				queue = queue[1:]
				for _, e := range edges {
					a, b := e.u, e.v
					if a < L || a > R || b < L || b > R {
						continue
					}
					if a == x {
						if c, ok := color[b]; !ok {
							color[b] = color[x] ^ 1
							queue = append(queue, b)
						} else if c == color[x] {
							return false
						}
					} else if b == x {
						if c, ok := color[a]; !ok {
							color[a] = color[x] ^ 1
							queue = append(queue, a)
						} else if c == color[x] {
							return false
						}
					}
				}
			}
		}
	}
	return true
}

func solveC(n int, edges []Edge, queries [][2]int) []int64 {
	res := make([]int64, len(queries))
	for qi, qr := range queries {
		L, R := qr[0], qr[1]
		var count int64
		for x := L; x <= R; x++ {
			for y := x; y <= R; y++ {
				if isBipartiteSub(n, edges, x, y) {
					count++
				}
			}
		}
		res[qi] = count
	}
	return res
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func distance(n int, edges []Edge, a, b int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{a}
	dist[a] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if v == b {
			return dist[v]
		}
		for _, u := range adj[v] {
			if dist[u] == -1 {
				dist[u] = dist[v] + 1
				q = append(q, u)
			}
		}
	}
	return -1
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 3
		edges := make([]Edge, 0)
		// build tree
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges = append(edges, Edge{p, i})
		}
		// maybe add extra edge forming odd cycle
		if rng.Intn(2) == 0 {
			for {
				u := rng.Intn(n) + 1
				v := rng.Intn(n) + 1
				if u == v {
					continue
				}
				d := distance(n, edges, u, v)
				if d != -1 && d%2 == 0 {
					edges = append(edges, Edge{u, v})
					break
				}
			}
		}
		m := len(edges)
		qn := rng.Intn(5) + 1
		queries := make([][2]int, qn)
		for i := 0; i < qn; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = [2]int{l, r}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
		fmt.Fprintf(&sb, "%d\n", qn)
		for _, qr := range queries {
			fmt.Fprintf(&sb, "%d %d\n", qr[0], qr[1])
		}
		input := sb.String()
		answers := solveC(n, edges, queries)
		var exp strings.Builder
		for i, a := range answers {
			if i > 0 {
				exp.WriteByte('\n')
			}
			fmt.Fprintf(&exp, "%d", a)
		}
		expected := exp.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
