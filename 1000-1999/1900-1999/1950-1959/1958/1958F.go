package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1_000_000_007

func modPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = int(int64(res) * int64(a) % int64(MOD))
		}
		a = int(int64(a) * int64(a) % int64(MOD))
		b >>= 1
	}
	return res
}

func combPrecompute(maxN int) ([]int, []int) {
	fac := make([]int, maxN+1)
	ifac := make([]int, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = int(int64(fac[i-1]) * int64(i) % int64(MOD))
	}
	ifac[maxN] = modPow(fac[maxN], MOD-2)
	for i := maxN; i > 0; i-- {
		ifac[i-1] = int(int64(ifac[i]) * int64(i) % int64(MOD))
	}
	return fac, ifac
}

func C(n, r int, fac, ifac []int) int {
	if r < 0 || r > n {
		return 0
	}
	return int(int64(fac[n]) * int64(ifac[r]) % int64(MOD) * int64(ifac[n-r]) % int64(MOD))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	maxF := 2 * n
	fac, ifac := combPrecompute(maxF)

	ans := make([]int, n+1)
	for s := 0; s <= 2*n-2; s++ {
		val := 0
		if s <= 2*n-4 && k >= 2 && k-2 <= 2*n-4-s {
			a := 0
			if s-(n-2) > a {
				a = s - (n - 2)
			}
			b := s
			if n-2 < b {
				b = n - 2
			}
			if b >= a {
				cnt := b - a + 1
				v := C(2*n-4-s, k-2, fac, ifac)
				val = (val + int(int64(cnt)*int64(v)%int64(MOD))) % MOD
			}
		}
		if s >= n-1 && s <= 2*n-3 && k >= 1 && k-1 <= 2*n-3-s {
			v := C(2*n-3-s, k-1, fac, ifac)
			val = (val + int((2*int64(v))%int64(MOD))) % MOD
		}
		if s < n-2 {
			ans[0] = (ans[0] + val) % MOD
		} else {
			idx := s - (n - 2)
			if idx <= n {
				ans[idx] = (ans[idx] + val) % MOD
			}
		}
	}

	for i := 0; i <= n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i]%MOD)
	}
	fmt.Fprintln(writer)
}
