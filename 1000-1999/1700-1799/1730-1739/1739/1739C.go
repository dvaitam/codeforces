package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	const N = 60
	var C [N + 1][N + 1]int64
	for i := 0; i <= N; i++ {
		C[i][0] = 1
		C[i][i] = 1
		for j := 1; j < i; j++ {
			C[i][j] = C[i-1][j-1] + C[i-1][j]
		}
	}

	var A, B [N + 1]int64
	A[2] = 1
	for n := 4; n <= N; n += 2 {
		A[n] = C[n-1][n/2-1] + B[n-2]
		B[n] = C[n][n/2] - A[n] - 1
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		fmt.Fprintf(out, "%d %d %d\n", A[n]%MOD, B[n]%MOD, 1)
	}
}
