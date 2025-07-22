package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Edge struct {
	to, rev, cap, cost int
}

func addEdge(g [][]Edge, u, v, cap, cost int) {
	g[u] = append(g[u], Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
	g[v] = append(g[v], Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

const oo = 100000

func solve(n, m int, a []int) string {
	v := make([]int, n+1)
	for i := 1; i <= n; i++ {
		j := i + 1
		for j <= n && a[j] != a[i] {
			j++
		}
		if j <= n {
			v[i] = j
		}
	}
	S := 0
	T := 3*n + 1
	tot := T + 1
	g := make([][]Edge, tot)
	addEdge(g, S, 1, m, 0)
	edIdx := make([]int, n+1)
	for i := 1; i <= n; i++ {
		edIdx[i] = len(g[i])
		addEdge(g, i, n+i, 1, bits.OnesCount(uint(a[i])))
		addEdge(g, n+i, 2*n+i, 1, -oo)
		addEdge(g, 2*n+i, T, 1, 0)
		if i < n {
			addEdge(g, 2*n+i, i+1, 1, 0)
			addEdge(g, i, i+1, oo, 0)
		}
	}
	for i := 1; i <= n; i++ {
		if v[i] > 0 {
			addEdge(g, 2*n+i, n+v[i], 1, 0)
		}
	}
	N := tot
	dist := make([]int, N)
	prevV := make([]int, N)
	prevE := make([]int, N)
	inQ := make([]bool, N)
	queue := make([]int, N)
	ans := 0
	for {
		for i := 0; i < N; i++ {
			dist[i] = 1 << 60
			inQ[i] = false
		}
		head, tail := 0, 0
		dist[S] = 0
		queue[tail] = S
		tail++
		inQ[S] = true
		for head < tail {
			u := queue[head]
			head++
			inQ[u] = false
			for ei, e := range g[u] {
				if e.cap > 0 && dist[u]+e.cost < dist[e.to] {
					dist[e.to] = dist[u] + e.cost
					prevV[e.to] = u
					prevE[e.to] = ei
					if !inQ[e.to] {
						inQ[e.to] = true
						queue[tail] = e.to
						tail++
					}
				}
			}
		}
		if dist[T] >= 0 {
			break
		}
		ans += dist[T]
		for vtx := T; vtx != S; {
			u := prevV[vtx]
			ei := prevE[vtx]
			g[u][ei].cap--
			rev := g[u][ei].rev
			g[vtx][rev].cap++
			vtx = u
		}
	}
	used := make([]bool, n+1)
	sel := 0
	for i := 1; i <= n; i++ {
		if g[i][edIdx[i]].cap == 0 {
			used[i] = true
			sel++
		}
	}
	m2 := 2*n - sel
	cost := ans % oo
	if cost < 0 {
		cost += oo
	}
	var out strings.Builder
	fmt.Fprintf(&out, "%d %d\n", m2, cost)
	now := make([]rune, n+2)
	bUsed := make([]bool, 256)
	for i := 1; i <= n; i++ {
		var w rune
		if now[i] == 0 {
			for c := 'a'; c <= 'z'; c++ {
				if !bUsed[c] {
					w = c
					break
				}
			}
			bUsed[w] = true
			fmt.Fprintf(&out, "%c=%d\n", w, a[i])
			now[i] = w
		} else {
			w = now[i]
		}
		fmt.Fprintf(&out, "print(%c)\n", w)
		if used[i] && v[i] > 0 {
			now[v[i]] = w
		} else {
			bUsed[w] = false
		}
	}
	return out.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) ([]int, int, int, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(3) + 1
	if m > n {
		m = n
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	return a, n, m, sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a, n, m, input := genCase(rng)
		exp := solve(n, m, a)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d mismatch:\nexpected:\n%s\ngot:\n%s", i+1, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
