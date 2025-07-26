package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	var s string
	fmt.Fscan(reader, &s)
	b := []byte(s)

	inv := make([]int64, n)
	for i := 1; i < n; i++ {
		inv[i] = modInv(int64(i))
	}

	ans := int64(1)
	for i := 2; i <= n-1; i++ {
		if b[i-1] == '?' {
			ans = ans * int64(i-1) % MOD
		}
	}
	zero := false
	if len(b) > 0 && b[0] == '?' {
		zero = true
	}
	if zero {
		fmt.Fprintln(writer, 0)
	} else {
		fmt.Fprintln(writer, ans)
	}

	for ; m > 0; m-- {
		var idx int
		var c string
		fmt.Fscan(reader, &idx, &c)
		idx--
		ch := c[0]
		if idx == 0 {
			if b[0] == '?' && ch != '?' {
				zero = false
			} else if b[0] != '?' && ch == '?' {
				zero = true
			}
			b[0] = ch
		} else {
			if b[idx] == '?' && ch != '?' {
				ans = ans * inv[idx-1] % MOD
			} else if b[idx] != '?' && ch == '?' {
				ans = ans * int64(idx-1) % MOD
			}
			b[idx] = ch
		}
		if zero {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}
