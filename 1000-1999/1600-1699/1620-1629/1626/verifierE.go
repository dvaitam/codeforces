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

func bfs(start int, g [][]int) []int {
	n := len(g)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func expected(n int, colors []int, edges [][2]int) string {
	g := make([][]int, n)
	blacks := []int{}
	for i, c := range colors {
		if c == 1 {
			blacks = append(blacks, i)
		}
	}
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	if len(blacks) >= 3 {
		ans := make([]byte, 0, 2*n-1)
		for i := 0; i < n; i++ {
			if i > 0 {
				ans = append(ans, ' ')
			}
			ans = append(ans, '1')
		}
		return string(ans)
	}
	if len(blacks) != 2 {
		return ""
	}
	b1 := blacks[0]
	b2 := blacks[1]
	d1 := bfs(b1, g)
	d2 := bfs(b2, g)
	dist := d1[b2]
	ans := make([]byte, 0, 2*n)
	for i := 0; i < n; i++ {
		t := (d1[i] + d2[i] - dist) / 2
		w := d1[i] - t
		ch := '0'
		if w <= 1 || dist-w <= 1 {
			ch = '1'
		}
		if i > 0 {
			ans = append(ans, ' ')
		}
		ans = append(ans, byte(ch))
	}
	return string(ans)
}

func runCase(bin string, n int, colors []int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(n))
	sb.WriteByte('\n')
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(c))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	input := sb.String()
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(n, colors, edges)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func randCase(rng *rand.Rand) (int, []int, [][2]int) {
	n := rng.Intn(10) + 2
	colors := make([]int, n)
	cnt := 0
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			colors[i] = 1
			cnt++
		} else {
			colors[i] = 0
		}
	}
	if cnt < 2 {
		colors[0] = 1
		colors[1] = 1
	}
	edges := randTree(rng, n)
	return n, colors, edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []struct {
		n int
		c []int
		e [][2]int
	}
	cases = append(cases, struct {
		n int
		c []int
		e [][2]int
	}{2, []int{1, 1}, [][2]int{{0, 1}}})
	cases = append(cases, struct {
		n int
		c []int
		e [][2]int
	}{3, []int{1, 0, 1}, [][2]int{{0, 1}, {1, 2}}})
	for i := 0; i < 100; i++ {
		n, c, e := randCase(rng)
		cases = append(cases, struct {
			n int
			c []int
			e [][2]int
		}{n, c, e})
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc.n, tc.c, tc.e); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
