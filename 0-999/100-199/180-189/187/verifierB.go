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

const INF int64 = 1 << 60

func floyd(dist [][]int64) {
	n := len(dist)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if v := dist[i][k] + dist[k][j]; v < dist[i][j] {
					dist[i][j] = v
				}
			}
		}
	}
}

func solveCase(input string) []int64 {
	in := strings.Split(strings.TrimSpace(input), "\n")
	// We'll parse manually using fmt
	reader := strings.NewReader(strings.Join(in, "\n"))
	var n, m, r int
	fmt.Fscan(reader, &n, &m, &r)
	distC := make([][][]int64, m)
	for c := 0; c < m; c++ {
		distC[c] = make([][]int64, n)
		for i := 0; i < n; i++ {
			distC[c][i] = make([]int64, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(reader, &distC[c][i][j])
			}
		}
		floyd(distC[c])
	}
	D0 := make([][]int64, n)
	for i := 0; i < n; i++ {
		D0[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			best := INF
			for c := 0; c < m; c++ {
				if distC[c][i][j] < best {
					best = distC[c][i][j]
				}
			}
			D0[i][j] = best
		}
	}
	answers := make([]int64, r)
	for idx := 0; idx < r; idx++ {
		var s, t, k int
		fmt.Fscan(reader, &s, &t, &k)
		s--
		t--
		seg := k + 1
		dpPrev := make([]int64, n)
		for i := range dpPrev {
			dpPrev[i] = INF
		}
		dpPrev[s] = 0
		for step := 0; step < seg; step++ {
			dpCurr := make([]int64, n)
			for i := range dpCurr {
				dpCurr[i] = INF
			}
			for u := 0; u < n; u++ {
				if dpPrev[u] == INF {
					continue
				}
				for v := 0; v < n; v++ {
					if val := dpPrev[u] + D0[u][v]; val < dpCurr[v] {
						dpCurr[v] = val
					}
				}
			}
			dpPrev = dpCurr
		}
		answers[idx] = dpPrev[t]
	}
	return answers
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 1
	r := rng.Intn(2) + 1
	distC := make([][][]int64, m)
	for c := 0; c < m; c++ {
		distC[c] = make([][]int64, n)
		for i := 0; i < n; i++ {
			distC[c][i] = make([]int64, n)
			for j := 0; j < n; j++ {
				if i == j {
					distC[c][i][j] = 0
				} else {
					distC[c][i][j] = int64(rng.Intn(9) + 1)
				}
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, r)
	for c := 0; c < m; c++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(distC[c][i][j]))
			}
			sb.WriteByte('\n')
		}
	}
	queries := make([][3]int, r)
	for i := 0; i < r; i++ {
		s := rng.Intn(n) + 1
		t := rng.Intn(n) + 1
		for t == s {
			t = rng.Intn(n) + 1
		}
		k := rng.Intn(2)
		queries[i] = [3]int{s, t, k}
		fmt.Fprintf(&sb, "%d %d %d\n", s, t, k)
	}
	ans := solveCase(sb.String())
	return sb.String(), ans
}

func runCase(bin string, in string, exp []int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outR := strings.Fields(strings.TrimSpace(out.String()))
	if len(outR) < len(exp) {
		return fmt.Errorf("expected %d lines, got %d", len(exp), len(outR))
	}
	for i, v := range exp {
		var got int64
		fmt.Sscan(outR[i], &got)
		if got != v {
			return fmt.Errorf("line %d: expected %d got %d", i+1, v, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
