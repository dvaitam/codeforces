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

type edge struct {
	x, y int
	w    int64
}

func bitsCount(x int) int {
	cnt := 0
	for x != 0 {
		x &= x - 1
		cnt++
	}
	return cnt
}

func shortestCycle(n int, edges []edge) int64 {
	deg := make([]int, n+1)
	var sumW int64
	adj := make([][]int, n+1)
	for _, e := range edges {
		sumW += e.w
		if e.x == e.y {
			deg[e.x] += 2
			if len(adj[e.x]) == 0 {
				adj[e.x] = append(adj[e.x], e.y)
			}
		} else {
			deg[e.x]++
			deg[e.y]++
			adj[e.x] = append(adj[e.x], e.y)
			adj[e.y] = append(adj[e.y], e.x)
		}
	}
	vis := make([]bool, n+1)
	q := []int{1}
	vis[1] = true
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				q = append(q, v)
			}
		}
	}
	for v := 1; v <= n; v++ {
		if deg[v] > 0 && !vis[v] {
			return -1
		}
	}
	const INF = int64(4e18)
	dist := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int64, n+1)
		for j := 1; j <= n; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = INF
			}
		}
	}
	for _, e := range edges {
		if e.w < dist[e.x][e.y] {
			dist[e.x][e.y] = e.w
			dist[e.y][e.x] = e.w
		}
	}
	for k := 1; k <= n; k++ {
		for i := 1; i <= n; i++ {
			if dist[i][k] == INF {
				continue
			}
			for j := 1; j <= n; j++ {
				nd := dist[i][k] + dist[k][j]
				if nd < dist[i][j] {
					dist[i][j] = nd
				}
			}
		}
	}
	odds := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if deg[i]%2 != 0 {
			odds = append(odds, i)
		}
	}
	k := len(odds)
	if k == 0 {
		return sumW
	}
	full := 1 << k
	dp := make([]int64, full)
	for mask := 1; mask < full; mask++ {
		dp[mask] = INF
	}
	dp[0] = 0
	for mask := 1; mask < full; mask++ {
		if bitsCount(mask)%2 != 0 {
			continue
		}
		var i int
		for i = 0; i < k; i++ {
			if mask&(1<<i) != 0 {
				break
			}
		}
		for j := i + 1; j < k; j++ {
			if mask&(1<<j) != 0 {
				m2 := mask ^ (1 << i) ^ (1 << j)
				cost := dp[m2] + dist[odds[i]][odds[j]]
				if cost < dp[mask] {
					dp[mask] = cost
				}
			}
		}
	}
	return sumW + dp[full-1]
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(10)
	edges := make([]edge, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		w := int64(rng.Intn(20) + 1)
		edges[i] = edge{x, y, w}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, w))
	}
	expected := fmt.Sprintf("%d", shortestCycle(n, edges))
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
