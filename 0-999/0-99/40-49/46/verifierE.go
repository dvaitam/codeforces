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

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func expected(mat [][]int64) int64 {
	n := len(mat)
	m := len(mat[0])
	ps := make([][]int64, n)
	for i := 0; i < n; i++ {
		ps[i] = make([]int64, m+1)
		for j := 1; j <= m; j++ {
			ps[i][j] = ps[i][j-1] + mat[i][j-1]
		}
	}
	const INF int64 = -1 << 60
	prev := make([]int64, m+2)
	curr := make([]int64, m+2)
	for j := 1; j <= m; j++ {
		prev[j] = ps[0][j]
	}
	for i := 2; i <= n; i++ {
		pref := make([]int64, m+2)
		suff := make([]int64, m+3)
		pref[0] = INF
		for j := 1; j <= m; j++ {
			pref[j] = max64(pref[j-1], prev[j])
		}
		suff[m+1] = INF
		for j := m; j >= 1; j-- {
			suff[j] = max64(suff[j+1], prev[j])
		}
		row := ps[i-1]
		if i%2 == 0 {
			for j := 1; j <= m; j++ {
				val := suff[j+1]
				if val <= INF/2 {
					curr[j] = INF
				} else {
					curr[j] = row[j] + val
				}
			}
		} else {
			for j := 1; j <= m; j++ {
				val := pref[j-1]
				if val <= INF/2 {
					curr[j] = INF
				} else {
					curr[j] = row[j] + val
				}
			}
		}
		prev, curr = curr, prev
	}
	ans := int64(-1 << 60)
	for j := 1; j <= m; j++ {
		if prev[j] > ans {
			ans = prev[j]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(4) + 2
	m := rng.Intn(4) + 2
	mat := make([][]int64, n)
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			mat[i][j] = int64(rng.Intn(21) - 10)
			fmt.Fprintf(&sb, "%d ", mat[i][j])
		}
		sb.WriteByte('\n')
	}
	expect := expected(mat)
	return sb.String(), expect
}

func runCase(exe, input string, expect int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
