package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD = 1000000007

func run(bin, input string) (string, error) {
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

func expected(n, k, m int) int64 {
	dp := make([]int64, n+2)
	dp2 := make([]int64, n+2)
	for M := 1; M <= n; M++ {
		dp[M] = 1
	}
	for i := 1; i < k; i++ {
		for idx := 1; idx <= n; idx++ {
			dp2[idx] = 0
		}
		for M := i; M <= n; M++ {
			v := dp[M]
			if v == 0 {
				continue
			}
			stay := int64(M - i)
			if stay > 0 {
				dp2[M] = (dp2[M] + v*stay) % MOD
			}
			end := M + m
			if end > n {
				end = n
			}
			for y := M + 1; y <= end; y++ {
				dp2[y] = (dp2[y] + v) % MOD
			}
		}
		dp, dp2 = dp2, dp
	}
	var ans int64
	for M := k; M <= n; M++ {
		ans = (ans + dp[M]) % MOD
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	if k > 12 {
		k = 12
	}
	m := rng.Intn(4) + 1
	input := fmt.Sprintf("%d %d %d\n", n, k, m)
	return input, expected(n, k, m)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	fixed := [][3]int{
		{3, 3, 1},
		{4, 2, 2},
	}
	idx := 0
	for ; idx < len(fixed); idx++ {
		n := fixed[idx][0]
		k := fixed[idx][1]
		m := fixed[idx][2]
		inp := fmt.Sprintf("%d %d %d\n", n, k, m)
		exp := strconv.FormatInt(expected(n, k, m), 10)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, exp, out, inp)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		inp, expVal := generateCase(rng)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strconv.FormatInt(expVal, 10) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", idx+1, expVal, out, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
