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

const mod = 1000000007

type caseD struct {
	m int
	X []int
	Y []int
}

func generateCase(rng *rand.Rand) caseD {
	m := rng.Intn(5) + 1
	X := make([]int, m)
	Y := make([]int, m)
	sum := 0
	for i := 0; i < m; i++ {
		X[i] = rng.Intn(3)
		sum += X[i]
	}
	if sum == 0 {
		idx := rng.Intn(m)
		X[idx] = 1
		sum = 1
	}
	for i := 0; i < m; i++ {
		Y[i] = X[i] + rng.Intn(3)
	}
	return caseD{m, X, Y}
}

func binom(n, k int, c [][]int) int {
	if k < 0 || k > n {
		return 0
	}
	if c[n][k] != 0 {
		return c[n][k]
	}
	if k == 0 || k == n {
		c[n][k] = 1
	} else {
		c[n][k] = (binom(n-1, k-1, c) + binom(n-1, k, c)) % mod
	}
	return c[n][k]
}

func solveCase(tc caseD) int {
	N := 0
	for _, v := range tc.X {
		N += v
	}
	c := make([][]int, N+1)
	for i := range c {
		c[i] = make([]int, N+1)
	}
	for i := 0; i <= N; i++ {
		c[i][0] = 1
		c[i][i] = 1
	}
	for i := 2; i <= N; i++ {
		for j := 1; j < i; j++ {
			c[i][j] = (c[i-1][j-1] + c[i-1][j]) % mod
		}
	}
	dpPrev := make([]int, N+1)
	dpCurr := make([]int, N+1)
	dpPrev[0] = 1
	for idx := 0; idx < tc.m; idx++ {
		xi := tc.X[idx]
		yi := tc.Y[idx]
		for i := range dpCurr {
			dpCurr[i] = 0
		}
		for open := 0; open <= N; open++ {
			v := dpPrev[open]
			if v == 0 {
				continue
			}
			t := open + xi
			maxK := yi
			if maxK > t {
				maxK = t
			}
			for k := 0; k <= maxK; k++ {
				ways := c[t][k]
				val := v * ways % mod
				dpCurr[t-k] = (dpCurr[t-k] + val) % mod
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}
	return dpPrev[0]
}

func runCase(bin string, tc caseD) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", tc.m)
	for i, v := range tc.X {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", v)
	}
	input.WriteByte('\n')
	for i, v := range tc.Y {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", v)
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := fmt.Sprintf("%d", solveCase(tc))
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
