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

type caseC struct {
	n int
	a [][]int
}

func genCaseC(rng *rand.Rand) caseC {
	n := rng.Intn(6) + 2 // 2..7
	a := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			a[i][j] = rng.Intn(21) - 10
		}
	}
	return caseC{n, a}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveC(tc caseC) int {
	n := tc.n
	a := tc.a
	const INF = -1 << 60
	dpPrev := make([][]int, n+1)
	dpCur := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dpPrev[i] = make([]int, n+1)
		dpCur[i] = make([]int, n+1)
		for j := 0; j <= n; j++ {
			dpPrev[i][j] = INF
			dpCur[i][j] = INF
		}
	}
	dpPrev[1][1] = a[1][1]
	for s := 3; s <= 2*n; s++ {
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				dpCur[i][j] = INF
			}
		}
		for i1 := max(1, s-n); i1 <= n && i1 < s; i1++ {
			j1 := s - i1
			if j1 < 1 || j1 > n {
				continue
			}
			for i2 := max(1, s-n); i2 <= n && i2 < s; i2++ {
				j2 := s - i2
				if j2 < 1 || j2 > n {
					continue
				}
				best := dpPrev[i1][i2]
				if i1 > 1 {
					best = max(best, dpPrev[i1-1][i2])
				}
				if i2 > 1 {
					best = max(best, dpPrev[i1][i2-1])
				}
				if i1 > 1 && i2 > 1 {
					best = max(best, dpPrev[i1-1][i2-1])
				}
				if best <= INF/2 {
					continue
				}
				val := best + a[i1][j1]
				if i1 != i2 {
					val += a[i2][j2]
				}
				dpCur[i1][i2] = val
			}
		}
		dpPrev, dpCur = dpCur, dpPrev
	}
	return dpPrev[n][n]
}

func runC(bin string, tc caseC) error {
	var sb strings.Builder
	fmt.Fprintln(&sb, tc.n)
	for i := 1; i <= tc.n; i++ {
		for j := 1; j <= tc.n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, tc.a[i][j])
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveC(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCaseC(rng)
		if err := runC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
