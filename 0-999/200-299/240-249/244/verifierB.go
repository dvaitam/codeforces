package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n int64) int64 {
	s := strconv.FormatInt(n, 10)
	L := len(s)
	digits := []byte(s)
	dp := make([][2][2][1024]int64, L+1)
	dp[0][1][0][0] = 1
	for pos := 0; pos < L; pos++ {
		for tight := 0; tight < 2; tight++ {
			for started := 0; started < 2; started++ {
				for mask := 0; mask < 1024; mask++ {
					cur := dp[pos][tight][started][mask]
					if cur == 0 {
						continue
					}
					limit := 9
					if tight == 1 {
						limit = int(digits[pos] - '0')
					}
					for d := 0; d <= limit; d++ {
						nt := 0
						if tight == 1 && d == limit {
							nt = 1
						}
						ns := started
						nm := mask
						if started == 0 {
							if d != 0 {
								ns = 1
								nm = 1 << d
							}
						} else {
							nm = mask | (1 << d)
						}
						if ns == 1 && bits.OnesCount(uint(nm)) > 2 {
							continue
						}
						dp[pos+1][nt][ns][nm] += cur
					}
				}
			}
		}
	}
	var res int64
	for tight := 0; tight < 2; tight++ {
		for mask := 0; mask < 1024; mask++ {
			if bits.OnesCount(uint(mask)) <= 2 {
				res += dp[L][tight][1][mask]
			}
		}
	}
	return res
}

func runCase(exe string, n int64) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(n)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edges := []int64{1, 2, 3, 4, 5, 9, 10, 11, 12, 20, 99, 100, 101, 110, 111, 112, 120, 1000, 123456789, 987654321, 1000000000}
	for i, n := range edges {
		if err := runCase(exe, n); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		if err := runCase(exe, n); err != nil {
			fmt.Fprintf(os.Stderr, "random case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
