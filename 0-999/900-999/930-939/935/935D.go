package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < int(n); i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < int(n); i++ {
		fmt.Fscan(in, &b[i])
	}

	invM := modPow(m, MOD-2)
	inv2 := modPow(2, MOD-2)

	pGreater := int64(0)
	pEqual := int64(1)

	for i := 0; i < int(n); i++ {
		ai := a[i]
		bi := b[i]
		if ai == 0 && bi == 0 {
			pGt := ((m - 1) % MOD) * inv2 % MOD * invM % MOD
			pEqualProb := invM
			pGreater = (pGreater + pEqual*pGt%MOD) % MOD
			pEqual = pEqual * pEqualProb % MOD
		} else if ai == 0 {
			pGt := ((m - bi) % MOD) * invM % MOD
			pEqualProb := invM
			pGreater = (pGreater + pEqual*pGt%MOD) % MOD
			pEqual = pEqual * pEqualProb % MOD
		} else if bi == 0 {
			pGt := ((ai - 1) % MOD) * invM % MOD
			pEqualProb := invM
			pGreater = (pGreater + pEqual*pGt%MOD) % MOD
			pEqual = pEqual * pEqualProb % MOD
		} else {
			if ai > bi {
				pGreater = (pGreater + pEqual) % MOD
				pEqual = 0
				break
			} else if ai < bi {
				pEqual = 0
				break
			}
		}
	}

	fmt.Fprintln(out, pGreater%MOD)
}
