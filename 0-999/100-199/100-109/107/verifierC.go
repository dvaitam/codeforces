package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

type solverC struct {
	n   int
	adj [][]int
	inE []int
}

func newSolverC(n int, edges []edge) *solverC {
	adj := make([][]int, n)
	inE := make([]int, n)
	for _, e := range edges {
		u, v := e.u-1, e.v-1
		if u < 0 || u >= n || v < 0 || v >= n {
			continue
		}
		adj[u] = append(adj[u], v)
		inE[v] |= 1 << u
	}
	return &solverC{n, adj, inE}
}

func (s *solverC) hasCycleUtil(v int, vis, rec []bool) bool {
	vis[v] = true
	rec[v] = true
	for _, to := range s.adj[v] {
		if !vis[to] {
			if s.hasCycleUtil(to, vis, rec) {
				return true
			}
		} else if rec[to] {
			return true
		}
	}
	rec[v] = false
	return false
}

func (s *solverC) hasCycle() bool {
	vis := make([]bool, s.n)
	rec := make([]bool, s.n)
	for i := 0; i < s.n; i++ {
		if !vis[i] {
			if s.hasCycleUtil(i, vis, rec) {
				return true
			}
		}
	}
	return false
}

func (s *solverC) count(given []int) int64 {
	size := 1 << s.n
	dp := make([]int64, size)
	dp[0] = 1
	for mask := 0; mask < size-1; mask++ {
		if dp[mask] == 0 {
			continue
		}
		c := bits.OnesCount(uint(mask))
		if given[c] == -1 {
			for j := 0; j < s.n; j++ {
				if mask&(1<<j) == 0 && mask&s.inE[j] == s.inE[j] {
					dp[mask|1<<j] += dp[mask]
				}
			}
		} else {
			j := given[c]
			if mask&s.inE[j] == s.inE[j] {
				dp[mask|1<<j] += dp[mask]
			}
		}
	}
	return dp[size-1]
}

func (s *solverC) solve(year int64) (string, bool) {
	year -= 2000
	if s.hasCycle() {
		return "", false
	}
	given := make([]int, s.n)
	for i := 0; i < s.n; i++ {
		given[i] = -1
	}
	used := make([]bool, s.n)
	res := make([]int, s.n)
	for i := 0; i < s.n; i++ {
		found := false
		for j := 0; j < s.n; j++ {
			if !used[j] {
				given[i] = j
				used[j] = true
				cnt := s.count(given)
				if cnt >= year {
					res[i] = j
					found = true
					break
				}
				year -= cnt
				given[i] = -1
				used[j] = false
			}
		}
		if !found {
			return "", false
		}
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	return sb.String(), true
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(n * n)
		edges := make([]edge, 0, m)
		exist := make(map[[2]int]bool)
		for len(edges) < m {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			if a == b {
				continue
			}
			key := [2]int{a, b}
			if exist[key] {
				continue
			}
			exist[key] = true
			edges = append(edges, edge{a, b})
		}
		year := int64(rng.Intn(30) + 2001)
		solver := newSolverC(n, edges)
		expect, ok := solver.solve(int64(year))
		if !ok {
			expect = "The times have changed\n"
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, year, len(edges)))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
