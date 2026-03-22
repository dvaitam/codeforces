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

var (
	n, k int
	w    []int
	adj  [][]int
	d    [][]int
	f    [][]int
	b    []int
	c    []int
)

func dfs(u, parent int) {
	for j := 1; j <= n; j++ {
		f[u][j] = w[d[u][j]] + k
	}
	for _, v := range adj[u] {
		if v == parent {
			continue
		}
		dfs(v, u)
		for j := 1; j <= n; j++ {
			a := f[v][j] - k
			bch := f[v][b[v]]
			if a < bch {
				f[u][j] += a
			} else {
				f[u][j] += bch
			}
		}
	}
	best := 1
	for j := 2; j <= n; j++ {
		if f[u][j] < f[u][best] {
			best = j
		}
	}
	b[u] = best
}

func prt(u, parent, z int) {
	c[u] = z
	for _, v := range adj[u] {
		if v == parent {
			continue
		}
		if f[v][z]-k < f[v][b[v]] {
			prt(v, u, z)
		} else {
			prt(v, u, b[v])
		}
	}
}

func solve(input string) string {
	parts := strings.Fields(input)
	idx := 0
	n = toInt(parts[idx])
	idx++
	k = toInt(parts[idx])
	idx++
	w = make([]int, n+1)
	for i := 1; i < n; i++ {
		w[i] = toInt(parts[idx])
		idx++
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		x := toInt(parts[idx])
		idx++
		y := toInt(parts[idx])
		idx++
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	// compute distances by BFS
	d = make([][]int, n+1)
	for u := 1; u <= n; u++ {
		dist := make([]int, n+1)
		for i := 1; i <= n; i++ {
			dist[i] = -1
		}
		queue := []int{u}
		dist[u] = 0
		for qi := 0; qi < len(queue); qi++ {
			v := queue[qi]
			for _, to := range adj[v] {
				if dist[to] < 0 {
					dist[to] = dist[v] + 1
					queue = append(queue, to)
				}
			}
		}
		row := make([]int, n+1)
		for i := 1; i <= n; i++ {
			row[i] = dist[i]
		}
		d[u] = row
	}
	f = make([][]int, n+1)
	for i := 1; i <= n; i++ {
		f[i] = make([]int, n+1)
	}
	b = make([]int, n+1)
	c = make([]int, n+1)
	dfs(1, 0)
	rootBest := b[1]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", f[1][rootBest]))
	prt(1, 0, rootBest)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	return sb.String()
}

func toInt(s string) int { var x int; fmt.Sscan(s, &x); return x }

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	k := rng.Intn(10) + 1
	weights := make([]int, n-1)
	for i := range weights {
		if i == 0 {
			weights[i] = rng.Intn(5)
		} else {
			weights[i] = weights[i-1] + rng.Intn(5)
		}
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range weights {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	if n > 1 {
		sb.WriteByte('\n')
	}
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func computeCost(input string, centers []int) (int, error) {
	parts := strings.Fields(input)
	idx := 0
	nn := toInt(parts[idx]); idx++
	kk := toInt(parts[idx]); idx++
	ww := make([]int, nn+1)
	for i := 1; i < nn; i++ {
		ww[i] = toInt(parts[idx]); idx++
	}
	adjLocal := make([][]int, nn+1)
	for i := 0; i < nn-1; i++ {
		x := toInt(parts[idx]); idx++
		y := toInt(parts[idx]); idx++
		adjLocal[x] = append(adjLocal[x], y)
		adjLocal[y] = append(adjLocal[y], x)
	}
	// BFS for distances
	dd := make([][]int, nn+1)
	for u := 1; u <= nn; u++ {
		dist := make([]int, nn+1)
		for i := 1; i <= nn; i++ { dist[i] = -1 }
		dist[u] = 0
		queue := []int{u}
		for qi := 0; qi < len(queue); qi++ {
			v := queue[qi]
			for _, to := range adjLocal[v] {
				if dist[to] < 0 {
					dist[to] = dist[v] + 1
					queue = append(queue, to)
				}
			}
		}
		dd[u] = dist
	}
	// compute cost
	// For each node u, cost = w[dist(u, centers[u])] + k if centers[u] != centers[parent[u]]
	// Actually: build tree rooted at 1, sum up costs
	parentLocal := make([]int, nn+1)
	childrenLocal := make([][]int, nn+1)
	visited := make([]bool, nn+1)
	visited[1] = true
	queue := []int{1}
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		for _, to := range adjLocal[v] {
			if !visited[to] {
				visited[to] = true
				parentLocal[to] = v
				childrenLocal[v] = append(childrenLocal[v], to)
				queue = append(queue, to)
			}
		}
	}
	total := 0
	for u := 1; u <= nn; u++ {
		total += ww[dd[u][centers[u]]]
		if u == 1 || centers[u] != centers[parentLocal[u]] {
			total += kk
		}
	}
	return total, nil
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	// Parse expected cost
	expParts := strings.Fields(exp)
	gotParts := strings.Fields(got)
	if len(expParts) == 0 || len(gotParts) == 0 {
		if got != exp {
			return fmt.Errorf("expected\n%s\nbut got\n%s", exp, got)
		}
		return nil
	}
	expCost := toInt(expParts[0])
	gotCost := toInt(gotParts[0])
	if expCost != gotCost {
		return fmt.Errorf("cost mismatch: expected %d got %d", expCost, gotCost)
	}
	// Parse candidate centers and verify cost
	inputFields := strings.Fields(strings.TrimRight(input, "\n"))
	nn := toInt(inputFields[0])
	if len(gotParts) != nn+1 {
		return fmt.Errorf("expected %d+1 tokens, got %d", nn, len(gotParts))
	}
	centers := make([]int, nn+1)
	for i := 1; i <= nn; i++ {
		centers[i] = toInt(gotParts[i])
		if centers[i] < 1 || centers[i] > nn {
			return fmt.Errorf("center %d for node %d out of range", centers[i], i)
		}
	}
	actualCost, err := computeCost(strings.TrimRight(input, "\n"), centers)
	if err != nil {
		return fmt.Errorf("cost verification error: %v", err)
	}
	if actualCost != gotCost {
		return fmt.Errorf("claimed cost %d but actual cost %d", gotCost, actualCost)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expected := solve(strings.TrimRight(input, "\n"))
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
