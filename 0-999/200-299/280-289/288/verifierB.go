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

const MOD int64 = 1000000007

func modpow(a, e int64) int64 {
	if e == 0 {
		return 1
	}
	if e < 0 {
		return modpow(modpow(a, -e), MOD-2)
	}
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func expected(n, k int) int64 {
	fact := make([]int64, k+1)
	invfact := make([]int64, k+1)
	fact[0] = 1
	for i := 1; i <= k; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	for i := 0; i <= k; i++ {
		invfact[i] = modpow(fact[i], MOD-2)
	}
	var sum int64
	for m := 1; m <= k; m++ {
		cyc := fact[k-1] * invfact[k-m] % MOD
		var trees int64
		if m < k {
			trees = int64(m) * modpow(int64(k), int64(k-m-1)) % MOD
		} else {
			trees = 1
		}
		sum = (sum + cyc*trees) % MOD
	}
	rem := modpow(int64(n-k), int64(n-k))
	ans := sum * rem % MOD
	return ans
}

func runCase(bin string, n, k int) error {
	input := fmt.Sprintf("%d %d\n", n, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("cannot parse output %q", gotStr)
	}
	exp := expected(n, k)
	if got != exp {
		return fmt.Errorf("expected %d got %d (n=%d k=%d)", exp, got, n, k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []struct{ n, k int }{
		{1, 1},
		{2, 2},
		{3, 1},
		{8, 3},
		{10, 5},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		n := rng.Intn(30) + 1
		k := rng.Intn(8) + 1
		if k > n {
			k = n
		}
		cases = append(cases, struct{ n, k int }{n, k})
	}
	for i, c := range cases {
		if err := runCase(bin, c.n, c.k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
