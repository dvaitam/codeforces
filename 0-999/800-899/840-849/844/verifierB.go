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

type testCaseB struct {
	n, m int
	grid [][]int
}

func expectedB(tc testCaseB) int64 {
	n, m := tc.n, tc.m
	pow2 := make([]int64, 51)
	pow2[0] = 1
	for i := 1; i <= 50; i++ {
		pow2[i] = pow2[i-1] * 2
	}
	ans := int64(0)
	for i := 0; i < n; i++ {
		cnt0, cnt1 := 0, 0
		for j := 0; j < m; j++ {
			if tc.grid[i][j] == 0 {
				cnt0++
			} else {
				cnt1++
			}
		}
		ans += pow2[cnt0] - 1
		ans += pow2[cnt1] - 1
	}
	for j := 0; j < m; j++ {
		cnt0, cnt1 := 0, 0
		for i := 0; i < n; i++ {
			if tc.grid[i][j] == 0 {
				cnt0++
			} else {
				cnt1++
			}
		}
		ans += pow2[cnt0] - 1
		ans += pow2[cnt1] - 1
	}
	ans -= int64(n * m)
	return ans
}

func runCase(bin string, tc testCaseB) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", tc.grid[i][j])
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
	exp := expectedB(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func genCase(rng *rand.Rand, n, m int) testCaseB {
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				grid[i][j] = 0
			} else {
				grid[i][j] = 1
			}
		}
	}
	return testCaseB{n: n, m: m, grid: grid}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseB, 0, 100)
	for n := 1; n <= 10; n++ {
		for m := 1; m <= 10; m++ {
			cases = append(cases, genCase(rng, n, m))
		}
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
