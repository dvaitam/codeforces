package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// Embedded correct solver for 1578L
type solverEdge struct {
	u, v int
	w    int64
}

func solveCase(input string) string {
	r := strings.NewReader(strings.TrimSpace(input))
	var n, m int
	fmt.Fscan(r, &n, &m)

	if n == 0 {
		return ""
	}

	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &c[i])
	}

	edges := make([]solverEdge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &edges[i].u, &edges[i].v, &edges[i].w)
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w > edges[j].w
	})

	dsu := make([]int, 2*n)
	for i := 1; i < 2*n; i++ {
		dsu[i] = i
	}
	var find func(i int) int
	find = func(i int) int {
		if dsu[i] == i {
			return i
		}
		dsu[i] = find(dsu[i])
		return dsu[i]
	}

	minI64 := func(a, b int64) int64 {
		if a < b {
			return a
		}
		return b
	}

	val := make([]int64, 2*n)
	sum_c := make([]int64, 2*n)
	for i := 1; i <= n; i++ {
		sum_c[i] = c[i]
	}

	children := make([][2]int, 2*n)

	idx := n
	for _, e := range edges {
		ru := find(e.u)
		rv := find(e.v)
		if ru != rv {
			idx++
			dsu[ru] = idx
			dsu[rv] = idx
			val[idx] = e.w
			sum_c[idx] = sum_c[ru] + sum_c[rv]
			children[idx] = [2]int{ru, rv}
		}
	}

	// If graph is not connected
	if idx < 2*n-1 {
		return "-1"
	}

	max_W := make([]int64, 2*n)
	max_W[idx] = math.MaxInt64

	for i := idx; i > n; i-- {
		u := children[i][0]
		v := children[i][1]
		max_W[u] = minI64(max_W[i], val[i]+sum_c[u])
		max_W[v] = minI64(max_W[i], val[i]+sum_c[v])
	}

	ans := int64(-1)
	total_c := sum_c[idx]
	for i := 1; i <= n; i++ {
		w_start := max_W[i] - total_c
		if w_start > ans {
			ans = w_start
		}
	}

	if ans <= 0 {
		return "-1"
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(r *rand.Rand) string {
	n := r.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := r.Intn(maxEdges) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", r.Intn(10)+1))
	}
	for i := 0; i < m; i++ {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		for v == u {
			v = r.Intn(n) + 1
		}
		w := r.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := solveCase(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
