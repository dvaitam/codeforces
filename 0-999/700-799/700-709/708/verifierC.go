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

func solveCase(n int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stk := []int{1}
	parent[1] = 0
	for len(stk) > 0 {
		v := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		order = append(order, v)
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			stk = append(stk, u)
		}
	}

	limit := n / 2
	sz := make([]int, n+1)
	sub := make([]int, n+1)

	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		s := 1
		best := 0
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			s += sz[v]
			if sub[v] > best {
				best = sub[v]
			}
		}
		sz[u] = s
		if s <= limit && s > best {
			best = s
		}
		sub[u] = best
	}

	up := make([]int, n+1)
	for _, u := range order {
		top1, top2 := 0, 0
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			val := sub[v]
			if val >= top1 {
				top2 = top1
				top1 = val
			} else if val > top2 {
				top2 = val
			}
		}
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			best := up[u]
			if sub[v] == top1 {
				if top2 > best {
					best = top2
				}
			} else {
				if top1 > best {
					best = top1
				}
			}
			comp := n - sz[v]
			if comp <= limit && comp > best {
				best = comp
			}
			up[v] = best
		}
	}

	var sb strings.Builder
	for u := 1; u <= n; u++ {
		mx := n - sz[u]
		best := up[u]
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			if sz[v] > mx {
				mx = sz[v]
				best = sub[v]
			}
		}
		if u > 1 {
			sb.WriteByte(' ')
		}
		if mx <= limit || mx-best <= limit {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

func runCase(bin string, n int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := solveCase(n, edges)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func randomTree(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(20) + 2
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return n, edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, edges := randomTree(rng)
		if err := runCase(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
