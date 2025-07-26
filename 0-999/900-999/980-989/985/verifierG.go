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

type pair struct{ u, v int }

func expectedG(n int, A, B, C uint64, edges []pair) uint64 {
	conflict := make([][]bool, n)
	for i := 0; i < n; i++ {
		conflict[i] = make([]bool, n)
	}
	for _, e := range edges {
		conflict[e.u][e.v] = true
		conflict[e.v][e.u] = true
	}
	var ans uint64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if conflict[i][j] {
				continue
			}
			for k := j + 1; k < n; k++ {
				if conflict[i][k] || conflict[j][k] {
					continue
				}
				ans += A*uint64(i) + B*uint64(j) + C*uint64(k)
			}
		}
	}
	return ans
}

func generateCaseG(rng *rand.Rand) (int, uint64, uint64, uint64, []pair) {
	n := rng.Intn(6) + 3
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	A := uint64(rng.Intn(10) + 1)
	B := uint64(rng.Intn(10) + 1)
	C := uint64(rng.Intn(10) + 1)
	used := make(map[pair]bool)
	edges := make([]pair, 0, m)
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		p := pair{u, v}
		if used[p] {
			continue
		}
		used[p] = true
		edges = append(edges, p)
	}
	return n, A, B, C, edges
}

func runCaseG(bin string, n int, A, B, C uint64, edges []pair) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	sb.WriteString(fmt.Sprintf("%d %d %d\n", A, B, C))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got uint64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := expectedG(n, A, B, C, edges)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, A, B, C, edges := generateCaseG(rng)
		if err := runCaseG(bin, n, A, B, C, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
