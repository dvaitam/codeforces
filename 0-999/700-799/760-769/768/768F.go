package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, e int64) int64 {
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

func modInv(a int64) int64 { return modPow(a, MOD-2) }

var fac, ifac []int64

func comb(n, k int) int64 {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var f, w, h int
	if _, err := fmt.Fscan(reader, &f, &w, &h); err != nil {
		return
	}
	if w == 0 {
		fmt.Println(1)
		return
	}
	if f == 0 {
		if w > h {
			fmt.Println(1)
		} else {
			fmt.Println(0)
		}
		return
	}
	limit := f + w + 5
	fac = make([]int64, limit)
	ifac = make([]int64, limit)
	fac[0] = 1
	for i := 1; i < limit; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[limit-1] = modInv(fac[limit-1])
	for i := limit - 1; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
	var total, liked int64
	for x := 1; x <= f; x++ {
		for _, dy := range []int{-1, 0, 1} {
			y := x + dy
			if y < 1 || y > w {
				continue
			}
			if abs(x-y) > 1 {
				continue
			}
			patterns := int64(1)
			if x == y {
				patterns = 2
			}
			waysF := comb(f-1, x-1)
			waysW := comb(w-1, y-1)
			total = (total + patterns*waysF%MOD*waysW) % MOD
			if w >= (h+1)*y {
				waysWL := comb(w-(h+1)*y+y-1, y-1)
				liked = (liked + patterns*waysF%MOD*waysWL) % MOD
			}
		}
	}
	invTotal := modInv(total)
	ans := liked % MOD * invTotal % MOD
	fmt.Println(ans)
}
