package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

const MOD int32 = 1000000007

func expected(s string) int32 {
	n := len(s)
	bs := []byte(" " + s)
	lcp := make([][]uint16, n+2)
	for i := range lcp {
		lcp[i] = make([]uint16, n+2)
	}
	for i := n; i >= 1; i-- {
		for j := n; j >= 1; j-- {
			if bs[i] == bs[j] {
				lcp[i][j] = lcp[i+1][j+1] + 1
			}
		}
	}
	dp := make([][]int32, n+2)
	pref := make([][]int32, n+2)
	for i := range dp {
		dp[i] = make([]int32, n+2)
		pref[i] = make([]int32, n+2)
	}
	dp[0][0] = 1
	pref[0][0] = 1
	for j := 1; j <= n; j++ {
		for i := 1; i <= j; i++ {
			if bs[i] == '0' {
				dp[i][j] = 0
				pref[i][j] = (pref[i-1][j] + dp[i][j]) % MOD
				continue
			}
			ans := pref[i-1][i-1]
			length := j - i + 1
			if k := i - length - 1; k >= 0 {
				val := pref[k][i-1]
				ans = (ans - val) % MOD
			}
			if k := i - length; k >= 1 {
				common := int(lcp[k][i])
				if common >= length || bs[k+common] > bs[i+common] {
					ans = (ans - dp[k][i-1]) % MOD
				}
			}
			if ans < 0 {
				ans += MOD
			}
			dp[i][j] = ans
			pref[i][j] = (pref[i-1][j] + dp[i][j]) % MOD
		}
	}
	return pref[n][n]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(4))
	digits := []byte("123456789")
	for t := 0; t < 100; t++ {
		n := rng.Intn(12) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(digits[rng.Intn(len(digits))])
		}
		s := sb.String()
		input := fmt.Sprintf("%d\n%s\n", n, s)
		want := expected(s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		got64, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil || int32(got64) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
