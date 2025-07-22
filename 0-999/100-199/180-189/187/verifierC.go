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

const INF int = 1 << 30

func bfs(n int, edges [][]int, src int) []int {
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := []int{src}
	dist[src] = 0
	for qi := 0; qi < len(q); qi++ {
		x := q[qi]
		for _, y := range edges[x] {
			if dist[y] == -1 {
				dist[y] = dist[x] + 1
				q = append(q, y)
			}
		}
	}
	return dist
}

func expectedAnswer(n, m, k int, vols []int, edgesU, edgesV []int, s, t int) int {
	edges := make([][]int, n+1)
	for i := 0; i < m; i++ {
		u := edgesU[i]
		v := edgesV[i]
		edges[u] = append(edges[u], v)
		edges[v] = append(edges[v], u)
	}
	distAll := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		distAll[i] = bfs(n, edges, i)
	}
	if distAll[s][t] < 0 {
		return -1
	}
	volIdx := make(map[int]int)
	for i, v := range vols {
		volIdx[v] = i
	}
	kVol := len(vols)
	distVolT := make([]int, kVol)
	for i, v := range vols {
		distVolT[i] = distAll[v][t]
	}
	startIdx := volIdx[s]
	maxQ := distAll[s][t]
	for q := 0; q <= maxQ; q++ {
		adj := make([][]int, kVol)
		for i := 0; i < kVol; i++ {
			for j := i + 1; j < kVol; j++ {
				if distAll[vols[i]][vols[j]] <= q {
					adj[i] = append(adj[i], j)
					adj[j] = append(adj[j], i)
				}
			}
		}
		vis := make([]bool, kVol)
		qarr := []int{startIdx}
		vis[startIdx] = true
		for qi := 0; qi < len(qarr); qi++ {
			u := qarr[qi]
			if distVolT[u] <= q {
				return q
			}
			for _, v := range adj[u] {
				if !vis[v] {
					vis[v] = true
					qarr = append(qarr, v)
				}
			}
		}
	}
	return maxQ
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(6) + 2
	m := rng.Intn(n*(n-1)/2-1) + n
	k := rng.Intn(n-1) + 1
	edgesU := make([]int, m)
	edgesV := make([]int, m)
	edgeSet := make(map[[2]int]struct{})
	for i := 0; i < m; {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := edgeSet[key]; ok {
			continue
		}
		edgeSet[key] = struct{}{}
		edgesU[i] = u
		edgesV[i] = v
		i++
	}
	vols := rng.Perm(n)[:k]
	for i := range vols {
		vols[i]++
	}
	s := vols[0]
	t := rng.Intn(n) + 1
	for t == s {
		t = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i, v := range vols {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", edgesU[i], edgesV[i])
	}
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	ans := expectedAnswer(n, m, k, vols, edgesU, edgesV, s, t)
	return sb.String(), ans
}

func runCase(bin string, input string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
