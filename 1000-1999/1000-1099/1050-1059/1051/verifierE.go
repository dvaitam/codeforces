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

const mod = 998244353

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func zfunc(s string) []int {
	n := len(s)
	z := make([]int, n)
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i <= r {
			z[i] = min(r-i+1, z[i-l])
		}
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
		if i+z[i]-1 > r {
			l = i
			r = i + z[i] - 1
		}
	}
	return z
}

func solveCase(a, L, R string) int64 {
	n := len(a)
	lLen := len(L)
	rLen := len(R)
	t1 := L + "#" + a
	z1 := zfunc(t1)
	t2 := R + "#" + a
	z2 := zfunc(t2)
	dp := make([]int64, n+1)
	sum := make([]int64, n+2)
	dp[n] = 1
	sum[n] = 1
	for i := n - 1; i >= 0; i-- {
		if a[i] == '0' {
			if lLen == 1 && L[0] == '0' {
				dp[i] = dp[i+1]
				sum[i] = (sum[i+1] + dp[i]) % mod
			} else {
				dp[i] = 0
				sum[i] = sum[i+1]
			}
			continue
		}
		if n-i < lLen {
			dp[i] = 0
			sum[i] = sum[i+1]
			continue
		}
		indL := lLen + 1 + i
		indR := rLen + 1 + i
		lf := i + lLen
		if z1[indL] != lLen && L[z1[indL]] > a[i+z1[indL]] {
			lf++
		}
		var rg int
		if n-i < rLen {
			rg = n
		} else if z2[indR] == rLen || R[z2[indR]] > a[i+z2[indR]] {
			rg = i + rLen
		} else {
			rg = i + rLen - 1
		}
		if rg >= lf {
			add := (sum[lf] - sum[rg+1] + mod) % mod
			dp[i] = add
			sum[i] = (sum[i+1] + add) % mod
		} else {
			dp[i] = 0
			sum[i] = sum[i+1]
		}
	}
	return dp[0] % mod
}

func genCase(rng *rand.Rand) (string, string, string) {
	n := rng.Intn(8) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('0' + rng.Intn(10)))
	}
	a := sb.String()
	L := fmt.Sprintf("%d", rng.Intn(100))
	R := fmt.Sprintf("%d", rng.Intn(100)+rng.Intn(100))
	if len(L) > len(R) || (len(L) == len(R) && L > R) {
		L, R = R, L
	}
	return a, L, R
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a, L, R := genCase(rng)
		input := fmt.Sprintf("%s\n%s\n%s\n", a, L, R)
		expect := solveCase(a, L, R)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(out, &got); err != nil || got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
