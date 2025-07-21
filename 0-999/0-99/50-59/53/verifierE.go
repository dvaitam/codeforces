package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func det(mat [][]int64) int64 {
	n := len(mat)
	if n == 0 {
		return 1
	}
	M := make([][]int64, n)
	for i := range mat {
		M[i] = make([]int64, n)
		copy(M[i], mat[i])
	}
	for k := 0; k < n; k++ {
		if M[k][k] == 0 {
			return 0
		}
		for i := k + 1; i < n; i++ {
			for j := k + 1; j < n; j++ {
				num := M[i][j]*M[k][k] - M[i][k]*M[k][j]
				den := int64(1)
				if k > 0 {
					den = M[k-1][k-1]
				}
				M[i][j] = num / den
			}
			M[i][k] = 0
		}
	}
	return M[n-1][n-1]
}

func countWays(n int, edges [][2]int, k int) int64 {
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
	}
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u][v] = true
		adj[v][u] = true
	}
	var total int64
	fullMask := 1 << n
	for mask := 0; mask < fullMask; mask++ {
		if bits.OnesCount(uint(mask)) != k {
			continue
		}
		rmask := (fullMask - 1) ^ mask
		prod := int64(1)
		valid := true
		for v := 0; v < n; v++ {
			if mask&(1<<v) != 0 {
				deg := 0
				for u := 0; u < n; u++ {
					if rmask&(1<<u) != 0 && adj[v][u] {
						deg++
					}
				}
				if deg == 0 {
					valid = false
					break
				}
				prod *= int64(deg)
			}
		}
		if !valid {
			continue
		}
		var rverts []int
		for v := 0; v < n; v++ {
			if rmask&(1<<v) != 0 {
				rverts = append(rverts, v)
			}
		}
		rlen := len(rverts)
		tcount := int64(1)
		if rlen > 1 {
			L := make([][]int64, rlen)
			for i := range L {
				L[i] = make([]int64, rlen)
			}
			for i := 0; i < rlen; i++ {
				for j := i + 1; j < rlen; j++ {
					u := rverts[i]
					v := rverts[j]
					if adj[u][v] {
						L[i][i]++
						L[j][j]++
						L[i][j]--
						L[j][i]--
					}
				}
			}
			sz := rlen - 1
			M := make([][]int64, sz)
			for i := 0; i < sz; i++ {
				M[i] = make([]int64, sz)
				for j := 0; j < sz; j++ {
					M[i][j] = L[i][j]
				}
			}
			tcount = det(M)
			if tcount == 0 {
				continue
			}
		}
		total += prod * tcount
	}
	return total
}

func generateGraph(rng *rand.Rand) (string, int64) {
	n := rng.Intn(5) + 3
	// build connected graph
	edges := make([][2]int, 0)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{p, i})
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float64() < 0.3 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	m := len(edges)
	k := rng.Intn(n-2) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	exp := countWays(n, edges, k)
	return sb.String(), exp
}

func runCase(exe, input string, expected int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateGraph(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
