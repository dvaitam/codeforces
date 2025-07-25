package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type edge struct{ u, v int }

func bfs(adj [][]int, start, goal []int) int {
	n := len(start)
	encode := func(arr []int) string {
		b := make([]byte, n)
		for i, v := range arr {
			b[i] = byte(v)
		}
		return string(b)
	}
	startKey := encode(start)
	goalKey := encode(goal)
	if startKey == goalKey {
		return 0
	}
	q := []string{startKey}
	dist := map[string]int{startKey: 0}
	states := map[string][]int{startKey: append([]int(nil), start...)}
	for len(q) > 0 {
		curKey := q[0]
		q = q[1:]
		arr := states[curKey]
		var zero int
		for i, v := range arr {
			if v == 0 {
				zero = i
				break
			}
		}
		for _, nb := range adj[zero] {
			next := make([]int, n)
			copy(next, arr)
			next[zero], next[nb] = next[nb], next[zero]
			key := encode(next)
			if _, ok := dist[key]; !ok {
				dist[key] = dist[curKey] + 1
				if key == goalKey {
					return dist[key]
				}
				states[key] = next
				q = append(q, key)
			}
		}
	}
	return -1
}

func solveF(n int, a, b []int, edges []edge) (int, int, int) {
	adj := make([][]int, n)
	exist := make(map[[2]int]bool)
	for _, e := range edges {
		u, v := e.u-1, e.v-1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		if u < v {
			exist[[2]int{u, v}] = true
		} else {
			exist[[2]int{v, u}] = true
		}
	}
	st := make([]int, n)
	copy(st, a)
	goal := make([]int, n)
	copy(goal, b)
	t0 := bfs(adj, st, goal)
	if t0 != -1 {
		return 0, 0, t0
	}
	bestT := -1
	bestU, bestV := 0, 0
	for u := 0; u < n; u++ {
		for v := u + 1; v < n; v++ {
			key := [2]int{u, v}
			if exist[key] {
				continue
			}
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			t := bfs(adj, st, goal)
			if t != -1 && (bestT == -1 || t < bestT) {
				bestT = t
				bestU = u + 1
				bestV = v + 1
			}
			adj[u] = adj[u][:len(adj[u])-1]
			adj[v] = adj[v][:len(adj[v])-1]
		}
	}
	if bestT == -1 {
		return -1, 0, 0
	}
	return bestU, bestV, bestT
}

func generateF(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 3
	a := rand.Perm(n)
	b := rand.Perm(n)
	edges := make([]edge, n-1)
	for i := 1; i < n; i++ {
		v := rng.Intn(i)
		edges[i-1] = edge{i + 1, v + 1}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range a {
		if v == n-1 {
			v = 0
		} else {
			v = v + 1
		}
		a[i] = v
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if v == n-1 {
			v = 0
		} else {
			v = v + 1
		}
		b[i] = v
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	u, v, t := solveF(n, a, b, edges)
	if u == -1 {
		return sb.String(), "-1"
	}
	if u == 0 {
		return sb.String(), fmt.Sprintf("0 %d", t)
	}
	if u > v {
		u, v = v, u
	}
	return sb.String(), fmt.Sprintf("%d %d %d", u, v, t)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(47))
	for i := 0; i < 100; i++ {
		input, exp := generateF(rng)
		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
