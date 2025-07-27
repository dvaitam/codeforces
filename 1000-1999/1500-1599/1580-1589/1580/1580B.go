package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	var p int64
	if _, err := fmt.Fscan(in, &n, &m, &k, &p); err != nil {
		return
	}

	// Precompute combinations C[n][k] modulo p
	C := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int64, n+1)
	}
	for i := 0; i <= n; i++ {
		C[i][0] = 1 % p
		C[i][i] = 1 % p
		for j := 1; j < i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % p
		}
	}

	fac := make([]int64, n+1)
	fac[0] = 1 % p
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % p
	}

	// dp[depth][size][count]
	dp := make([][][]int64, m)
	for d := 0; d < m; d++ {
		dp[d] = make([][]int64, n+1)
		for s := 0; s <= n; s++ {
			dp[d][s] = make([]int64, n+1)
		}
	}

	dp[0][0][0] = 1 % p
	for size := 1; size <= n; size++ {
		dp[0][size][1] = fac[size]
	}

	for depth := 1; depth < m; depth++ {
		for size := 0; size <= n; size++ {
			if size == 0 {
				dp[depth][size][0] = 1 % p
				continue
			}
			for L := 0; L < size; L++ {
				R := size - 1 - L
				comb := C[size-1][L]
				for a := 0; a <= L; a++ {
					if dp[depth-1][L][a] == 0 {
						continue
					}
					for b := 0; b <= R && a+b <= size; b++ {
						if dp[depth-1][R][b] == 0 {
							continue
						}
						t := a + b
						val := comb * dp[depth-1][L][a] % p
						val = val * dp[depth-1][R][b] % p
						dp[depth][size][t] += val
						if dp[depth][size][t] >= p {
							dp[depth][size][t] %= p
						}
					}
				}
			}
		}
	}

	fmt.Fprintln(out, dp[m-1][n][k]%p)
}
