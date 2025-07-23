package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func matMul(a, b [][]int64, d int) [][]int64 {
	res := make([][]int64, d)
	for i := 0; i < d; i++ {
		res[i] = make([]int64, d)
		for k := 0; k < d; k++ {
			if a[i][k] == 0 {
				continue
			}
			aik := a[i][k]
			for j := 0; j < d; j++ {
				res[i][j] = (res[i][j] + aik*b[k][j]) % mod
			}
		}
	}
	return res
}

func matPow(m [][]int64, d int, exp int64) [][]int64 {
	// identity matrix
	res := make([][]int64, d)
	for i := 0; i < d; i++ {
		res[i] = make([]int64, d)
		res[i][i] = 1
	}
	for exp > 0 {
		if exp&1 == 1 {
			res = matMul(res, m, d)
		}
		m = matMul(m, m, d)
		exp >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var b int64
	var k, x int
	if _, err := fmt.Fscan(in, &n, &b, &k, &x); err != nil {
		return
	}
	cnt := make([]int64, 10)
	for i := 0; i < n; i++ {
		var d int
		fmt.Fscan(in, &d)
		cnt[d]++
	}
	// build transition matrix of size x
	M := make([][]int64, x)
	for i := 0; i < x; i++ {
		M[i] = make([]int64, x)
		for d := 1; d <= 9; d++ {
			if cnt[d] == 0 {
				continue
			}
			to := (i*10 + d) % x
			M[i][to] = (M[i][to] + cnt[d]) % mod
		}
	}
	Mp := matPow(M, x, b)
	ans := Mp[0][k] % mod
	fmt.Println(ans)
}
