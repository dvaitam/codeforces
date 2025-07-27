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

const mod = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func solveD(n int) int64 {
	dp := make([]int64, n+1)
	sum0 := make([]int64, n+1)
	sum1 := make([]int64, n+1)
	dp[0] = 1
	sum0[0] = 1
	for i := 1; i <= n; i++ {
		if (i-1)&1 == 0 {
			dp[i] = sum0[i-1]
		} else {
			dp[i] = sum1[i-1]
		}
		if dp[i] >= mod {
			dp[i] %= mod
		}
		if i&1 == 0 {
			sum0[i] = (sum0[i-1] + dp[i]) % mod
			sum1[i] = sum1[i-1]
		} else {
			sum1[i] = (sum1[i-1] + dp[i]) % mod
			sum0[i] = sum0[i-1]
		}
	}
	total := modPow(2, int64(n))
	invTotal := modPow(total, mod-2)
	return dp[n] * invTotal % mod
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("run error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	for i := 0; i < t; i++ {
		n := rand.Intn(50) + 1
		exp := solveD(n)
		input := fmt.Sprintf("%d\n", n)
		out, err := runBinary(binary, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		vStr := strings.TrimSpace(out)
		v, err := strconv.ParseInt(vStr, 10, 64)
		if err != nil || v != exp {
			fmt.Printf("test %d failed: n=%d expected=%d got=%s\n", i+1, n, exp, vStr)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
