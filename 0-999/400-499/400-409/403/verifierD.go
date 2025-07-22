package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const MOD = 1000000007
const MAXN = 1000
const Kmax = 44

func add(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}

func mul(a, b int) int {
	return int((int64(a) * int64(b)) % MOD)
}

var f [MAXN + 1][Kmax + 1]int

func precompute() {
	var D [Kmax + 1][MAXN + 1]int
	D[0][0] = 1
	for d := 0; d <= MAXN; d++ {
		for k := Kmax; k >= 1; k-- {
			for s := d; s <= MAXN; s++ {
				if D[k-1][s-d] != 0 {
					D[k][s] = add(D[k][s], D[k-1][s-d])
				}
			}
		}
	}
	fact := make([]int, MAXN+1)
	invfact := make([]int, MAXN+1)
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = mul(fact[i-1], i)
	}
	invfact[MAXN] = modInv(fact[MAXN])
	for i := MAXN; i > 0; i-- {
		invfact[i-1] = mul(invfact[i], i)
	}
	comb := func(n, k int) int {
		if n < 0 || k < 0 || n < k {
			return 0
		}
		return mul(fact[n], mul(invfact[k], invfact[n-k]))
	}
	factk := make([]int, Kmax+1)
	factk[0] = 1
	for i := 1; i <= Kmax; i++ {
		factk[i] = mul(factk[i-1], i)
	}
	for k := 1; k <= Kmax; k++ {
		minS := k * (k - 1) / 2
		for s := minS; s <= MAXN; s++ {
			cnt := D[k][s]
			if cnt == 0 {
				continue
			}
			coef := mul(factk[k], cnt)
			for n := k + s; n <= MAXN; n++ {
				f[n][k] = add(f[n][k], mul(coef, comb(n-s, k)))
			}
		}
	}
}

func modInv(a int) int { return modPow(a, MOD-2) }

func modPow(a, e int) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 != 0 {
			res = mul(res, base)
		}
		base = mul(base, base)
		e >>= 1
	}
	return res
}

func solve(n, k int) int {
	if k < 0 || k > Kmax || k > n {
		return 0
	}
	return f[n][k]
}

func runCase(bin string, n, k int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%d %d\n", n, k))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got int
	fmt.Sscan(outStr, &got)
	exp := solve(n, k)
	if got != exp {
		return fmt.Errorf("expected %d got %s", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	precompute()
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, k int
		fmt.Sscan(line, &n, &k)
		if err := runCase(bin, n, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
