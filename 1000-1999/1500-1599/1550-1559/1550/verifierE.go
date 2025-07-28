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

func possible(s []byte, n, k, L int) bool {
	if L == 0 {
		return true
	}
	if k*L > n {
		return false
	}
	const INF = int(1e9)
	next := make([][]int, k)
	prefix := make([]int, n+1)
	for c := 0; c < k; c++ {
		prefix[0] = 0
		target := byte('a' + c)
		for i := 0; i < n; i++ {
			if s[i] != '?' && s[i] != target {
				prefix[i+1] = prefix[i] + 1
			} else {
				prefix[i+1] = prefix[i]
			}
		}
		nxt := make([]int, n+1)
		nxt[n] = INF
		for i := n - 1; i >= 0; i-- {
			nxt[i] = nxt[i+1]
			if i+L <= n && prefix[i+L]-prefix[i] == 0 {
				nxt[i] = i
			}
		}
		next[c] = nxt
	}
	dp := make([]int, 1<<k)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for mask := 0; mask < (1 << k); mask++ {
		pos := dp[mask]
		if pos > n {
			continue
		}
		for c := 0; c < k; c++ {
			if mask&(1<<c) != 0 {
				continue
			}
			start := next[c][pos]
			if start == INF {
				continue
			}
			end := start + L
			if end <= n && end < dp[mask|(1<<c)] {
				dp[mask|(1<<c)] = end
			}
		}
	}
	return dp[(1<<k)-1] <= n
}

func computeExpected(n, k int, s string) int {
	bytes := []byte(s)
	low, high := 0, n
	ans := 0
	for low <= high {
		mid := (low + high) / 2
		if possible(bytes, n, k, mid) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return ans
}

type testCase struct {
	n        int
	k        int
	s        string
	expected int
}

func generateCase(rng *rand.Rand) testCase {
	k := rng.Intn(3) + 1 // 1..3 letters
	n := rng.Intn(10) + k
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		r := rng.Intn(k + 1)
		if r == k {
			bytes[i] = '?'
		} else {
			bytes[i] = byte('a' + r)
		}
	}
	s := string(bytes)
	exp := computeExpected(n, k, s)
	return testCase{n: n, k: k, s: s, expected: exp}
}

func (tc testCase) input() string {
	return fmt.Sprintf("%d %d\n%s\n", tc.n, tc.k, tc.s)
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
	var got int
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
