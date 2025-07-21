package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func uniqueInt64(a []int64) []int64 {
	j := 0
	for i := 0; i < len(a); i++ {
		if i == 0 || a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func solveCase(n int, dist [][]int64, g1, g2, s1, s2 int) int64 {
	const INF = int64(1e18)
	// Floyd-Warshall
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == INF {
				continue
			}
			for j := 0; j < n; j++ {
				if dist[k][j] == INF {
					continue
				}
				nd := dist[i][k] + dist[k][j]
				if nd < dist[i][j] {
					dist[i][j] = nd
				}
			}
		}
	}
	dvals := make([]int64, 0, n*(n-1))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				dvals = append(dvals, dist[i][j])
			}
		}
	}
	sort.Slice(dvals, func(i, j int) bool { return dvals[i] < dvals[j] })
	dvals = uniqueInt64(dvals)
	D := len(dvals)
	A := make([][]int, D)
	for p := 0; p < D; p++ {
		A[p] = make([]int, n)
		for i := 0; i < n; i++ {
			cnt := 0
			for j := 0; j < n; j++ {
				if i != j && dist[i][j] <= dvals[p] {
					cnt++
				}
			}
			A[p][i] = cnt
		}
	}
	ans := int64(0)
	maxG := g2
	maxS := s2
	for p := 0; p < D; p++ {
		for q := p + 1; q < D; q++ {
			a := make([]int64, n)
			b := make([]int64, n)
			c := make([]int64, n)
			for i := 0; i < n; i++ {
				a[i] = int64(A[p][i])
				b[i] = int64(A[q][i] - A[p][i])
				c[i] = int64((n - 1) - A[q][i])
			}
			dp := make([][]int64, maxG+1)
			for i := range dp {
				dp[i] = make([]int64, maxS+1)
			}
			dp[0][0] = 1
			for i := 0; i < n; i++ {
				ndp := make([][]int64, maxG+1)
				for ii := range ndp {
					ndp[ii] = make([]int64, maxS+1)
				}
				for gi := 0; gi <= maxG; gi++ {
					for si := 0; si <= maxS; si++ {
						v := dp[gi][si]
						if v == 0 {
							continue
						}
						if gi < maxG && a[i] > 0 {
							ndp[gi+1][si] += v * a[i]
						}
						if si < maxS && b[i] > 0 {
							ndp[gi][si+1] += v * b[i]
						}
						if c[i] > 0 {
							ndp[gi][si] += v * c[i]
						}
					}
				}
				dp = ndp
			}
			for gi := g1; gi <= g2; gi++ {
				if gi > maxG {
					break
				}
				for si := s1; si <= s2; si++ {
					if si > maxS {
						break
					}
					ans += dp[gi][si]
				}
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 3
	dist := make([][]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = int64(rng.Intn(20) + 1)
			}
		}
	}
	g2 := rng.Intn(n)
	g1 := rng.Intn(g2 + 1)
	s2 := rng.Intn(n - g2)
	s1 := rng.Intn(s2 + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, n*(n-1)/2))
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			w := dist[i][j]
			sb.WriteString(fmt.Sprintf("%d %d %d\n", i+1, j+1, w))
			dist[j][i] = w
		}
	}
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", g1, g2, s1, s2))
	ans := solveCase(n, dist, g1, g2, s1, s2)
	return sb.String(), fmt.Sprintf("%d", ans)
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
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
