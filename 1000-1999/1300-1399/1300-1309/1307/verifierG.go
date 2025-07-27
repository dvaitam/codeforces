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

type Edge struct {
	to, rev, cap, cost int
}

func addEdge(g [][]Edge, u, v, cap, cost int) {
	g[u] = append(g[u], Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
	g[v] = append(g[v], Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(input string) string {
	in := strings.NewReader(input)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}
	g := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		addEdge(g, u-1, v-1, 1, c)
	}
	var q int
	fmt.Fscan(in, &q)
	x := make([]int64, q)
	ans := make([]float64, q)
	for i := 0; i < q; i++ {
		var xi int64
		fmt.Fscan(in, &xi)
		x[i] = xi
		ans[i] = 1e18
	}
	const INF int64 = 1 << 60
	sumF := 0
	var sumC int64
	dist := make([]int64, n)
	prevv := make([]int, n)
	preve := make([]int, n)
	inq := make([]bool, n)
	qv := make([]int, n+5)
	for {
		for i := 0; i < n; i++ {
			dist[i] = INF
			inq[i] = false
		}
		dist[0] = 0
		head, tail := 0, 0
		qv[tail] = 0
		tail++
		inq[0] = true
		for head < tail {
			v := qv[head]
			head++
			inq[v] = false
			for ei, e := range g[v] {
				if e.cap > 0 && dist[e.to] > dist[v]+int64(e.cost) {
					dist[e.to] = dist[v] + int64(e.cost)
					prevv[e.to] = v
					preve[e.to] = ei
					if !inq[e.to] {
						qv[tail] = e.to
						tail++
						inq[e.to] = true
					}
				}
			}
		}
		if dist[n-1] == INF {
			break
		}
		sumC += dist[n-1]
		sumF++
		v := n - 1
		for v != 0 {
			u := prevv[v]
			ei := preve[v]
			g[u][ei].cap--
			rev := g[u][ei].rev
			g[v][rev].cap++
			v = u
		}
		for i := 0; i < q; i++ {
			cur := float64(sumC+x[i]) / float64(sumF)
			if cur < ans[i] {
				ans[i] = cur
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%.12f\n", ans[i]))
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges-n+1) + n - 1
	edges := make([][3]int, m)
	used := make(map[[2]int]bool)
	for i := range edges {
		for {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if used[[2]int{u, v}] {
				continue
			}
			used[[2]int{u, v}] = true
			edges[i][0] = u
			edges[i][1] = v
			edges[i][2] = rng.Intn(10) + 1
			break
		}
	}
	q := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		fmt.Fprintf(&sb, "%d\n", rng.Intn(20))
	}
	in := sb.String()
	return in, solve(in)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
