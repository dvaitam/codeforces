package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Pair struct {
	dist int
	id   int
}

func comparePair(a, b Pair) int {
	if a.dist != b.dist {
		return a.dist - b.dist
	}
	return a.id - b.id
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// solve is the embedded correct solver for 38H, matching the CF-accepted solution.
func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, mEdges int
	fmt.Fscan(in, &n, &mEdges)

	dist := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = 1e7
			}
		}
	}

	for i := 0; i < mEdges; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		dist[u][v] = c
		dist[v][u] = c
	}

	var g1, g2, s1, s2 int
	fmt.Fscan(in, &g1, &g2, &s1, &s2)

	for k := 1; k <= n; k++ {
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	P := make([][]Pair, n+1)
	minP := make([]Pair, n+1)
	maxP := make([]Pair, n+1)

	for i := 1; i <= n; i++ {
		P[i] = make([]Pair, 0, n-1)
		for j := 1; j <= n; j++ {
			if i != j {
				P[i] = append(P[i], Pair{dist[i][j], i})
			}
		}
		mn := P[i][0]
		mx := P[i][0]
		for _, p := range P[i] {
			if comparePair(p, mn) < 0 {
				mn = p
			}
			if comparePair(p, mx) > 0 {
				mx = p
			}
		}
		minP[i] = mn
		maxP[i] = mx
	}

	ans := int64(0)
	var dp [55][55]int64

	for u := 1; u <= n; u++ {
		for v := 1; v <= n; v++ {
			if u == v {
				continue
			}
			if comparePair(minP[u], maxP[v]) >= 0 {
				continue
			}

			possible := true
			mask := make([]int, n+1)
			for i := 1; i <= n; i++ {
				if i == u {
					mask[i] = 1
				} else if i == v {
					mask[i] = 4
				} else {
					m := 0
					if comparePair(minP[i], minP[u]) < 0 {
						m |= 1
					}
					if comparePair(maxP[i], maxP[v]) > 0 {
						m |= 4
					}
					for _, p := range P[i] {
						if comparePair(minP[u], p) < 0 && comparePair(p, maxP[v]) < 0 {
							m |= 2
							break
						}
					}
					mask[i] = m
					if m == 0 {
						possible = false
						break
					}
				}
			}
			if !possible {
				continue
			}

			for j := 0; j <= g2; j++ {
				for k := 0; k <= s2; k++ {
					dp[j][k] = 0
				}
			}
			dp[0][0] = 1

			for i := 1; i <= n; i++ {
				m := mask[i]
				for j := minInt(i-1, g2); j >= 0; j-- {
					for k := minInt(i-1-j, s2); k >= 0; k-- {
						val := dp[j][k]
						if val == 0 {
							continue
						}
						dp[j][k] = 0
						if m&4 != 0 {
							dp[j][k] += val
						}
						if m&1 != 0 && j+1 <= g2 {
							dp[j+1][k] += val
						}
						if m&2 != 0 && k+1 <= s2 {
							dp[j][k+1] += val
						}
					}
				}
			}

			for j := g1; j <= g2; j++ {
				for k := s1; k <= s2; k++ {
					ans += dp[j][k]
				}
			}
		}
	}

	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 3
	edges := make([]struct{ u, v, w int }, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			w := rng.Intn(20) + 1
			edges = append(edges, struct{ u, v, w int }{i, j, w})
		}
	}
	g2 := rng.Intn(n)
	g1 := rng.Intn(g2 + 1)
	s2 := rng.Intn(n - g2)
	s1 := rng.Intn(s2 + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", g1, g2, s1, s2))
	return sb.String()
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		expected := solve(input)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
