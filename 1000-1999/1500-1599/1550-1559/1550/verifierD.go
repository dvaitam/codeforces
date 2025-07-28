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

type testCase struct {
	n        int
	l        int
	r        int
	expected int64
}

func brute(n, l, r int) int64 {
	arr := make([]int, n)
	maxF := -1
	var count int64
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == n {
			for i := 0; i < n; i++ {
				if arr[i] == i+1 {
					return
				}
			}
			f := 0
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					if arr[i]+arr[j] == (i+1)+(j+1) {
						f++
					}
				}
			}
			if f > maxF {
				maxF = f
				count = 1
			} else if f == maxF {
				count++
			}
			return
		}
		for v := l; v <= r; v++ {
			arr[pos] = v
			dfs(pos + 1)
		}
	}
	dfs(0)
	return count % mod
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 2     // 2..5
	width := rng.Intn(3) + 2 // 2..4 values
	l := rng.Intn(7) - 3
	r := l + width - 1
	exp := brute(n, l, r)
	return testCase{n: n, l: l, r: r, expected: exp}
}

func (tc testCase) input() string {
	return fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.l, tc.r)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input())
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
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
