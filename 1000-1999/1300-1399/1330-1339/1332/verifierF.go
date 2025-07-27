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

const MOD int64 = 998244353

type edge struct{ u, v int }

func solveCase(n int, edges []edge) int64 {
	m := len(edges)
	// adjacency matrix
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
	}
	for _, e := range edges {
		adj[e.u][e.v] = true
		adj[e.v][e.u] = true
	}
	var sum int64
	// iterate over non-empty subsets of edges
	for mask := 1; mask < (1 << m); mask++ {
		involvedMask := 0
		for i, e := range edges {
			if mask>>i&1 == 1 {
				involvedMask |= 1 << e.u
				involvedMask |= 1 << e.v
			}
		}
		// get list of vertices
		vertIndices := []int{}
		for v := 0; v < n; v++ {
			if involvedMask>>v&1 == 1 {
				vertIndices = append(vertIndices, v)
			}
		}
		k := len(vertIndices)
		var cnt int64
		for sub := 0; sub < (1 << k); sub++ {
			ok := true
			for i := 0; i < k && ok; i++ {
				if sub>>i&1 == 1 {
					vi := vertIndices[i]
					for j := i + 1; j < k; j++ {
						if sub>>j&1 == 1 {
							vj := vertIndices[j]
							if adj[vi][vj] {
								ok = false
								break
							}
						}
					}
				}
			}
			if ok {
				cnt++
			}
		}
		sum = (sum + cnt) % MOD
	}
	return sum
}

func generateCase(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(5) + 2 // 2..6
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p - 1, i - 1})
	}
	return n, edges
}

func runCase(bin string, n int, edges []edge) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		input.WriteString(fmt.Sprintf("%d %d\n", e.u+1, e.v+1))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := solveCase(n, edges)
	outStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(outStr, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got%MOD != expected%MOD {
		return fmt.Errorf("expected %d got %d", expected%MOD, got%MOD)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, edges := generateCase(rng)
		if err := runCase(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
