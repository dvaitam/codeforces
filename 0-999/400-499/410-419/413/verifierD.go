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

type testCase struct {
	n, k int
	seq  []int
}

func expected(tc testCase) int64 {
	mod := int64(1000000007)
	n := tc.n
	k := tc.k
	seq := tc.seq
	maxMask := 1 << (k - 1)
	suff := make([]int64, n)
	if n > 0 {
		suff[n-1] = 1
		for i := n - 2; i >= 0; i-- {
			mult := int64(1)
			if seq[i+1] == 0 {
				mult = 2
			}
			suff[i] = suff[i+1] * mult % mod
		}
	}
	maskNext := make([][]int, 3)
	willWin := make([][]bool, 3)
	for h := 1; h <= 2; h++ {
		maskNext[h] = make([]int, maxMask)
		willWin[h] = make([]bool, maxMask)
		for mask := 0; mask < maxMask; mask++ {
			t := h
			cur := mask
			for {
				if t == k {
					willWin[h][mask] = true
					break
				}
				bit := 1 << (t - 1)
				if cur&bit == 0 {
					cur |= bit
					maskNext[h][mask] = cur
					break
				}
				cur &^= bit
				t++
			}
		}
	}
	dp := make([]int64, maxMask)
	dp[0] = 1
	var ans int64
	for i, v := range seq {
		newDp := make([]int64, maxMask)
		hs := []int{}
		if v == 0 {
			hs = []int{1, 2}
		} else if v == 2 {
			hs = []int{1}
		} else if v == 4 {
			hs = []int{2}
		}
		for mask := 0; mask < maxMask; mask++ {
			ways := dp[mask]
			if ways == 0 {
				continue
			}
			for _, h := range hs {
				if willWin[h][mask] {
					ans = (ans + ways*suff[i]) % mod
				} else {
					nm := maskNext[h][mask]
					newDp[nm] = (newDp[nm] + ways) % mod
				}
			}
		}
		dp = newDp
	}
	return ans % mod
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(tc.seq[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	var got int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(15) + 1
	k := rng.Intn(5) + 3
	seq := make([]int, n)
	for i := 0; i < n; i++ {
		v := rng.Intn(3)
		if v == 0 {
			seq[i] = 0
		} else if v == 1 {
			seq[i] = 2
		} else {
			seq[i] = 4
		}
	}
	return testCase{n: n, k: k, seq: seq}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	tests = append(tests, testCase{n: 1, k: 3, seq: []int{2}})
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
