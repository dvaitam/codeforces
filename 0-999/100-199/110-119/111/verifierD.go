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

var c [1005][1005]int64
var f [1005]int64
var fz [2005]int64
var fm [1005]int64

func modExp(base, exp int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func getC(n int) {
	for i := 0; i <= n; i++ {
		c[i][0] = 1
		for j := 1; j <= i; j++ {
			c[i][j] = c[i-1][j] + c[i-1][j-1]
			if c[i][j] >= mod {
				c[i][j] -= mod
			}
		}
	}
}

func solveD(n, m, k int) int64 {
	if m == 1 {
		return modExp(int64(k), int64(n))
	}
	getC(n)
	for i := 1; i <= n; i++ {
		f[i] = modExp(int64(i), int64(n))
		for j := 1; j < i; j++ {
			f[i] = (f[i] - f[j]*c[i][j]%mod + mod) % mod
		}
	}
	fz[0], fm[0] = 1, 1
	limit := 2 * n
	if limit > k {
		limit = k
	}
	for i := 1; i <= limit; i++ {
		fz[i] = fz[i-1] * int64(k-i+1) % mod
	}
	for i := 1; i <= n; i++ {
		fm[i] = fm[i-1] * modExp(int64(i), mod-2) % mod
	}
	var ans int64
	for i := 0; i <= n && i <= k; i++ {
		temp := modExp(int64(i), int64(m-2)*int64(n))
		for j := 0; j <= n-i && i+2*j <= k; j++ {
			idx := i + 2*j
			term := temp * fz[idx] % mod
			term = term * fm[i] % mod
			term = term * fm[j] % mod
			term = term * fm[j] % mod
			val := f[i+j]
			term = term * val % mod
			term = term * val % mod
			ans = (ans + term) % mod
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	k := rng.Intn(20) + 1
	input := fmt.Sprintf("%d %d %d\n", n, m, k)
	out := fmt.Sprintf("%d\n", solveD(n, m, k))
	return input, out
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
