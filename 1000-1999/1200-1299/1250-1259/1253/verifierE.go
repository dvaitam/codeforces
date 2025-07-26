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

func expectedE(n, m int, ants [][2]int) int {
	type Ant struct{ x, s int }
	arr := make([]Ant, n)
	for i, v := range ants {
		arr[i] = Ant{v[0], v[1]}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].x < arr[j].x })
	const INF = int(1e9)
	dp := make([]int, m+1)
	for i := 1; i <= m; i++ {
		dp[i] = INF
	}
	for _, a := range arr {
		newdp := make([]int, m+1)
		copy(newdp, dp)
		best := make([]int, m+1)
		for i := 0; i <= m; i++ {
			best[i] = INF
		}
		for pos := 0; pos < m; pos++ {
			if dp[pos] == INF {
				continue
			}
			need := pos + 1
			d := 0
			if need < a.x-a.s {
				d = a.x - a.s - need
			}
			r0 := a.x + a.s + d
			if r0 > m {
				r0 = m
			}
			val := dp[pos] + d - r0
			if val < best[r0] {
				best[r0] = val
			}
		}
		pref := INF
		for r := 0; r <= m; r++ {
			if best[r] < pref {
				pref = best[r]
			}
			if pref == INF {
				continue
			}
			if cand := pref + r; cand < newdp[r] {
				newdp[r] = cand
			}
		}
		dp = newdp
	}
	return dp[m]
}

func generateCaseE(rng *rand.Rand) (int, int, [][2]int) {
	n := rng.Intn(4) + 1
	m := rng.Intn(10) + 1
	ants := make([][2]int, n)
	for i := 0; i < n; i++ {
		ants[i][0] = rng.Intn(m) + 1
		ants[i][1] = rng.Intn(5)
	}
	return n, m, ants
}

func runCaseE(bin string, n, m int, ants [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, v := range ants {
		sb.WriteString(fmt.Sprintf("%d %d\n", v[0], v[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expect := expectedE(n, m, ants)
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, ants := generateCaseE(rng)
		if err := runCaseE(bin, n, m, ants); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
