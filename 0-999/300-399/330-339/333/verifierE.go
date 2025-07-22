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

type edge struct {
	d    int
	u, v int
}

func solve(n int, xs, ys []int) float64 {
	tot := n * (n - 1) / 2
	edges := make([]edge, 0, tot)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			dx := xs[i] - xs[j]
			dy := ys[i] - ys[j]
			d2 := dx*dx + dy*dy
			edges = append(edges, edge{d: d2, u: i, v: j})
		}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].d > edges[j].d })
	g := (n + 63) / 64
	bits := make([][]uint64, n)
	for i := range bits {
		bits[i] = make([]uint64, g)
	}
	for _, e := range edges {
		u, v := e.u, e.v
		for k := 0; k < g; k++ {
			if bits[u][k]&bits[v][k] != 0 {
				return math.Sqrt(float64(e.d)) / 2.0
			}
		}
		bits[u][v>>6] |= 1 << (uint(v) & 63)
		bits[v][u>>6] |= 1 << (uint(u) & 63)
	}
	return 0.0
}

func runCase(bin string, n int, xs, ys []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solve(n, xs, ys)
	if math.Abs(got-exp) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", exp, got)
	}
	return nil
}

func genCase(rng *rand.Rand) (int, []int, []int) {
	n := rng.Intn(6) + 3
	xs := make([]int, n)
	ys := make([]int, n)
	used := make(map[[2]int]struct{})
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(21) - 10
			y := rng.Intn(21) - 10
			p := [2]int{x, y}
			if _, ok := used[p]; !ok {
				used[p] = struct{}{}
				xs[i] = x
				ys[i] = y
				break
			}
		}
	}
	return n, xs, ys
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct {
		n      int
		xs, ys []int
	}
	xs1 := []int{0, 1, 2}
	ys1 := []int{0, 0, 0}
	cases = append(cases, struct {
		n      int
		xs, ys []int
	}{3, xs1, ys1})
	for i := 0; i < 100; i++ {
		n, xs, ys := genCase(rng)
		cases = append(cases, struct {
			n      int
			xs, ys []int
		}{n, xs, ys})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.n, tc.xs, tc.ys); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
