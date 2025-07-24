package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	var k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	origK := k
	limitK := k
	if limitK > int(n) {
		limitK = int(n)
	}
	inv := make([]int64, origK+2)
	inv[1] = 1
	for i := 2; i <= origK; i++ {
		inv[i] = MOD - MOD/int64(i)*inv[int(MOD%int64(i))]%MOD
	}
	maxPair := origK
	if int64(maxPair) > n/2 {
		maxPair = int(n / 2)
	}
	pair := make([]int64, maxPair+1)
	pair[0] = 1
	for x := 1; x <= maxPair; x++ {
		a1 := (n - int64(2*x) + 2) % MOD
		a2 := (n - int64(2*x) + 1) % MOD
		d1 := int64(x)
		d2 := (n - int64(x) + 1) % MOD
		pair[x] = pair[x-1] * a1 % MOD
		pair[x] = pair[x] * a2 % MOD
		pair[x] = pair[x] * modPow(d1%MOD, MOD-2) % MOD
		pair[x] = pair[x] * modPow(d2, MOD-2) % MOD
	}
	ans := make([]int64, origK+1)
	for x := 0; x <= maxPair; x++ {
		base := pair[x]
		maxR := origK - x
		if int64(maxR) > n-int64(2*x) {
			maxR = int(n - int64(2*x))
		}
		singles := int64(1)
		if x <= origK {
			ans[x] = (ans[x] + base) % MOD
		}
		for r := 1; r <= maxR; r++ {
			singles = singles * ((n - int64(2*x) - int64(r) + 1) % MOD) % MOD
			singles = singles * modPow(int64(r)%MOD, MOD-2) % MOD
			ans[x+r] = (ans[x+r] + base*singles) % MOD
		}
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for i := 1; i <= origK; i++ {
		if i > 1 {
			fmt.Fprint(w, " ")
		}
		if i <= limitK {
			fmt.Fprint(w, ans[i]%MOD)
		} else {
			fmt.Fprint(w, 0)
		}
	}
	fmt.Fprintln(w)
}
