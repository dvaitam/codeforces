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

const MOD int64 = 998244353

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func expected(n, a0, x, y, k, M int64) int64 {
	L := int64(1)
	for i := int64(1); i <= k; i++ {
		L = L * i / gcd(L, i)
	}
	invN := modInv(n % MOD)
	mulStay := (MOD + 1 - invN) % MOD
	dp := make([]int64, L)
	next := make([]int64, L)
	for step := k; step >= 1; step-- {
		s := int64(step)
		for r := int64(0); r < L; r++ {
			m := r % s
			val := (invN*(r+dp[r-m]) + mulStay*dp[r]) % MOD
			next[r] = val
		}
		dp, next = next, dp
	}
	constPart := (L % MOD) * (k % MOD) % MOD * invN % MOD
	ans := int64(0)
	val := a0
	for i := int64(0); i < n; i++ {
		q := val / L
		r := val % L
		ans += dp[r]
		ans += (q % MOD) * constPart % MOD
		ans %= MOD
		val = (val*x + y) % M
	}
	ans = ans * modPow(n%MOD, k) % MOD
	return ans
}

func runCase(bin string, n, a0, x, y, k, M int64) error {
	input := fmt.Sprintf("%d %d %d %d %d %d\n", n, a0, x, y, k, M)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := fmt.Sprint(expected(n, a0, x, y, k, M))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randCase(rng *rand.Rand) (int64, int64, int64, int64, int64, int64) {
	n := int64(rng.Intn(5) + 1)
	a0 := int64(rng.Intn(10))
	x := int64(rng.Intn(10) + 1)
	y := int64(rng.Intn(10))
	k := int64(rng.Intn(4) + 1)
	M := int64(rng.Intn(20) + 1)
	return n, a0, x, y, k, M
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][6]int64
	cases = append(cases, [6]int64{1, 0, 1, 0, 1, 2})
	cases = append(cases, [6]int64{2, 3, 2, 1, 2, 5})
	for i := 0; i < 100; i++ {
		n, a0, x, y, k, M := randCase(rng)
		cases = append(cases, [6]int64{n, a0, x, y, k, M})
	}
	for idx, c := range cases {
		if err := runCase(bin, c[0], c[1], c[2], c[3], c[4], c[5]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
