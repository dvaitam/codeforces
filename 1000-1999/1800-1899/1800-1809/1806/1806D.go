package main

import (
	"bufio"
	"io"
	"os"
)

const MOD int64 = 998244353
const MAXN = 500000 + 5

func modAdd(a, b int64) int64 {
	x := a + b
	if x >= MOD {
		x -= MOD
	}
	return x
}
func modMul(a, b int64) int64 { return (a * b) % MOD }

func main() {
	data, _ := io.ReadAll(os.Stdin)
	idx := 0
	nextInt := func() int {
		n := len(data)
		for idx < n && (data[idx] <= ' ') {
			idx++
		}
		sign := 1
		if idx < n && data[idx] == '-' {
			sign = -1
			idx++
		}
		x := 0
		for idx < n && data[idx] >= '0' && data[idx] <= '9' {
			x = x*10 + int(data[idx]-'0')
			idx++
		}
		return x * sign
	}

	inv := make([]int64, MAXN)
	fact := make([]int64, MAXN)
	inv[1] = 1
	for i := 2; i < MAXN; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[int(MOD%int64(i))]%MOD
	}
	fact[0] = 1
	for i := 1; i < MAXN; i++ {
		fact[i] = modMul(fact[i-1], int64(i))
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	writeInt := func(x int64) {
		if x == 0 {
			out.WriteByte('0')
			return
		}
		if x < 0 {
			out.WriteByte('-')
			x = -x
		}
		var buf [20]byte
		i := len(buf)
		for x > 0 {
			i--
			buf[i] = byte('0' + x%10)
			x /= 10
		}
		out.Write(buf[i:])
	}

	t := nextInt()
	for ; t > 0; t-- {
		n := nextInt()
		a := make([]int, n)
		for i := 1; i <= n-1; i++ {
			a[i] = nextInt()
		}
		pref := int64(1)
		sum := int64(0)
		for k := 1; k <= n-1; k++ {
			if a[k] == 1 {
				pref = modMul(pref, modMul(int64(k-1), inv[k]))
			}
			if a[k] == 0 {
				add := modMul(inv[k], pref)
				sum = modAdd(sum, add)
			}
			ans := modMul(fact[k], sum)
			writeInt(ans)
			if k < n-1 {
				out.WriteByte(' ')
			}
		}
		if n-1 > 0 {
			out.WriteByte('\n')
		}
	}
}
