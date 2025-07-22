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

type DSU struct {
	p []int
}

func newDSU(n int) *DSU {
	d := &DSU{p: make([]int, n+1)}
	for i := range d.p {
		d.p[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(x, y int) {
	rx := d.find(x)
	ry := d.find(y)
	if rx != ry {
		d.p[ry] = rx
	}
}

type edge struct {
	u, v, w int
}

func expected(n, m, k int, c []int, edges []edge) string {
	typeOf := make([]int, n+1)
	idx := 1
	for i := 1; i <= k; i++ {
		cnt := c[i-1]
		for j := 0; j < cnt; j++ {
			typeOf[idx] = i
			idx++
		}
	}
	dsu := newDSU(n)
	for _, e := range edges {
		if e.w == 0 {
			dsu.union(e.u, e.v)
		}
	}
	idx = 1
	for i := 1; i <= k; i++ {
		cnt := c[i-1]
		if cnt > 0 {
			root := dsu.find(idx)
			for j := 0; j < cnt; j++ {
				if dsu.find(idx+j) != root {
					return "No"
				}
			}
		}
		idx += cnt
	}
	const INF = int(1e9)
	dist := make([][]int, k)
	for i := 0; i < k; i++ {
		dist[i] = make([]int, k)
		for j := 0; j < k; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = INF
			}
		}
	}
	for _, e := range edges {
		if e.w == 0 {
			continue
		}
		ti := typeOf[e.u] - 1
		tj := typeOf[e.v] - 1
		if ti != tj {
			if e.w < dist[ti][tj] {
				dist[ti][tj] = e.w
				dist[tj][ti] = e.w
			}
		}
	}
	for p := 0; p < k; p++ {
		for i := 0; i < k; i++ {
			if dist[i][p] == INF {
				continue
			}
			for j := 0; j < k; j++ {
				if dist[p][j] == INF {
					continue
				}
				nd := dist[i][p] + dist[p][j]
				if nd < dist[i][j] {
					dist[i][j] = nd
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("Yes\n")
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			if dist[i][j] >= INF {
				sb.WriteString("-1")
			} else {
				sb.WriteString(fmt.Sprintf("%d", dist[i][j]))
			}
			if j+1 < k {
				sb.WriteByte(' ')
			}
		}
		if i+1 < k {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		k := rng.Intn(3) + 1
		c := make([]int, k)
		n := 0
		for i := 0; i < k; i++ {
			c[i] = rng.Intn(3) + 1
			n += c[i]
		}
		m := rng.Intn(4*n + 1)
		var edges []edge
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < k; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", c[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			w := rng.Intn(5)
			edges = append(edges, edge{u, v, w})
			sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
		}
		input := sb.String()
		expectedOut := expected(n, m, k, c, edges)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expectedOut) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", tcase+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
