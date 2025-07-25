package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int64 = 998244353

func applyN(v []int64) []int64 {
	n := len(v)
	res := make([]int64, n)
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = (prefix[i-1] + v[i-1]) % MOD
	}
	for i := 1; i <= n; i++ {
		lb := i & -i
		res[i-1] = (prefix[i-1] - prefix[i-lb]) % MOD
		if res[i-1] < 0 {
			res[i-1] += MOD
		}
	}
	return res
}

func powMod(a, e int64) int64 {
	r := int64(1)
	for e > 0 {
		if e&1 == 1 {
			r = r * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return r
}

func modInv(a int64) int64 { return powMod(a, MOD-2) }

func solveCase(n int, k int64, b []int64) []int64 {
	r := bits.Len(uint(n))
	powers := make([][]int64, r)
	powers[0] = append([]int64(nil), b...)
	for i := 1; i < r; i++ {
		powers[i] = applyN(powers[i-1])
	}
	coeff := make([]int64, r)
	coeff[0] = 1
	comb := int64(1)
	for i := 1; i < r; i++ {
		comb = comb * (k + int64(i-1)) % MOD
		comb = comb * modInv(int64(i)) % MOD
		val := comb
		if i%2 == 1 {
			val = (MOD - val) % MOD
		}
		coeff[i] = val
	}
	ans := make([]int64, n)
	for i := 0; i < r; i++ {
		c := coeff[i]
		if c == 0 {
			continue
		}
		for j := 0; j < n; j++ {
			ans[j] = (ans[j] + c*powers[i][j]) % MOD
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		a := solveCase(n, k, b)
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, a[i])
		}
		writer.WriteByte('\n')
	}
}
