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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(n, k int, s, t string) int {
	t0 := t[0]
	t1 := t[1]
	if t0 == t1 {
		cnt := 0
		for i := 0; i < n; i++ {
			if s[i] == t0 {
				cnt++
			}
		}
		cnt = cnt + min(k, n-cnt)
		return cnt * (cnt - 1) / 2
	}
	const INF = -1_000_000_000
	dp := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]int, n+1)
		for j := 0; j <= n; j++ {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0
	for pos := 0; pos < n; pos++ {
		c := s[pos]
		ndp := make([][]int, k+1)
		for i := 0; i <= k; i++ {
			ndp[i] = make([]int, n+1)
			for j := 0; j <= n; j++ {
				ndp[i][j] = INF
			}
		}
		for used := 0; used <= k; used++ {
			for cnt0 := 0; cnt0 <= n; cnt0++ {
				cur := dp[used][cnt0]
				if cur < 0 {
					continue
				}
				newUsed := used
				newCnt0 := cnt0
				add := 0
				if c == t0 {
					newCnt0 = cnt0 + 1
				} else if c == t1 {
					add = cnt0
				}
				if cur+add > ndp[newUsed][newCnt0] {
					ndp[newUsed][newCnt0] = cur + add
				}
				if used < k {
					newUsed2 := used + 1
					newCnt02 := cnt0 + 1
					if cur > ndp[newUsed2][newCnt02] {
						ndp[newUsed2][newCnt02] = cur
					}
					add1 := cnt0
					if cur+add1 > ndp[newUsed2][cnt0] {
						ndp[newUsed2][cnt0] = cur + add1
					}
				}
			}
		}
		dp = ndp
	}
	best := 0
	for used := 0; used <= k; used++ {
		for cnt0 := 0; cnt0 <= n; cnt0++ {
			if dp[used][cnt0] > best {
				best = dp[used][cnt0]
			}
		}
	}
	return best
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rng.Intn(26) + 'a')
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 2
	k := rng.Intn(n + 1)
	s := randString(rng, n)
	t := randString(rng, 2)
	ans := solve(n, k, s, t)
	input := fmt.Sprintf("%d %d\n%s\n%s\n", n, k, s, t)
	expected := fmt.Sprintf("%d", ans)
	return input, expected
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
