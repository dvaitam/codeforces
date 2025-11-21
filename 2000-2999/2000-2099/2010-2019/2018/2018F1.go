package main

import (
	"bufio"
	"fmt"
	"os"
)

func powMod(a, e, mod int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var mod int64
		fmt.Fscan(in, &n, &mod)

		atLeast := make([][]int64, n+2)
		for i := 0; i <= n+1; i++ {
			atLeast[i] = make([]int64, n+2)
		}

		for L := 1; L <= n; L++ {
			for R := L; R <= n; R++ {
				prod := int64(1)
				valid := true
				for i := 1; i <= n; i++ {
					need := 1
					if i < L {
						need = 1 + (L - i)
					} else if i > R {
						need = 1 + (i - R)
					}
					if need > n {
						prod = 0
						valid = false
						break
					}
					choices := int64(n - need + 1)
					prod = prod * choices % mod
				}
				if !valid {
					prod = 0
				}
				atLeast[L][R] = prod
			}
		}

		exact := make([][]int64, n+2)
		for i := 0; i <= n+1; i++ {
			exact[i] = make([]int64, n+2)
		}

		res := make([]int64, n+1)
		sumNonZero := int64(0)

		for L := 1; L <= n; L++ {
			for R := L; R <= n; R++ {
				val := atLeast[L][R]
				if L > 1 {
					val = (val - atLeast[L-1][R]) % mod
				}
				if R < n {
					val = (val - atLeast[L][R+1]) % mod
				}
				if L > 1 && R < n {
					val = (val + atLeast[L-1][R+1]) % mod
				}
				if val < 0 {
					val += mod
				}
				exact[L][R] = val
				k := R - L + 1
				res[k] = (res[k] + val) % mod
				sumNonZero = (sumNonZero + val) % mod
			}
		}

		total := powMod(int64(n)%mod, int64(n), mod)
		res[0] = (total - sumNonZero) % mod
		if res[0] < 0 {
			res[0] += mod
		}

		for k := 0; k <= n; k++ {
			if k > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[k])
		}
		fmt.Fprintln(out)
	}
}
