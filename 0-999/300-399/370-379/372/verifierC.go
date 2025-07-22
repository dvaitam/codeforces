package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(n, m, d int, a, b, tVals []int) int64 {
	dp := make([]int64, n)
	dp2 := make([]int64, n)
	for x := 0; x < n; x++ {
		dp[x] = int64(b[0] - absInt(a[0]-x))
	}
	for i := 1; i < m; i++ {
		dist := (tVals[i] - tVals[i-1]) * d
		if dist >= n {
			var best int64 = math.MinInt64
			for x := 0; x < n; x++ {
				if dp[x] > best {
					best = dp[x]
				}
			}
			for x := 0; x < n; x++ {
				dp2[x] = best + int64(b[i]-absInt(a[i]-x))
			}
		} else {
			head, tail := 0, 0
			deque := make([]int, n)
			r := -1
			for x := 0; x < n; x++ {
				newR := x + dist
				if newR >= n {
					newR = n - 1
				}
				for r < newR {
					r++
					for tail > head && dp[deque[tail-1]] <= dp[r] {
						tail--
					}
					deque[tail] = r
					tail++
				}
				l := x - dist
				if l < 0 {
					l = 0
				}
				for head < tail && deque[head] < l {
					head++
				}
				best := dp[deque[head]]
				dp2[x] = best + int64(b[i]-absInt(a[i]-x))
			}
		}
		dp, dp2 = dp2, dp
	}
	var ans int64 = math.MinInt64
	for x := 0; x < n; x++ {
		if dp[x] > ans {
			ans = dp[x]
		}
	}
	return ans
}

func runCase(bin string, n, m, d int, a, b, tVals []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, d))
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a[i]+1, b[i], tVals[i]))
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
	expected := solveCase(n, m, d, a, b, tVals)
	var got int64
	fmt.Sscan(strings.TrimSpace(out.String()), &got)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	for i := 0; i < tests; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(3) + 1
		d := rng.Intn(3) + 1
		a := make([]int, m)
		b := make([]int, m)
		tVals := make([]int, m)
		curTime := 0
		for j := 0; j < m; j++ {
			a[j] = rng.Intn(n)
			b[j] = rng.Intn(10)
			curTime += rng.Intn(5) + 1
			tVals[j] = curTime
		}
		if err := runCase(bin, n, m, d, a, b, tVals); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
