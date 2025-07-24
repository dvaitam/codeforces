package main

import (
	"bufio"
	"fmt"
	"os"
)

func modPow(a, e, mod int64) int64 {
	res := int64(1 % mod)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInverse(a, mod int64) int64 {
	return modPow(a, mod-2, mod)
}

func multiplicativeOrder(a, mod int64) int64 {
	a %= mod
	cur := int64(1)
	for i := int64(1); ; i++ {
		cur = cur * a % mod
		if cur == 1 {
			return i
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M, a, Q int64
	if _, err := fmt.Fscan(in, &N, &M, &a, &Q); err != nil {
		return
	}

	phi := multiplicativeOrder(a, Q)

	limit := N - 1
	if M < limit {
		limit = M
	}
	prefix := make([]int64, limit+1)
	c := int64(1)
	prefix[0] = 1 % phi
	for i := int64(1); i <= limit; i++ {
		c = c * (M - i + 1) % phi
		c = c * modInverse(i%phi, phi) % phi
		prefix[i] = (prefix[i-1] + c) % phi
	}
	pow2M := modPow(2, M, phi)

	for i := int64(1); i <= N; i++ {
		k := N - i
		var exp int64
		if k > M {
			exp = pow2M
		} else {
			exp = prefix[k]
		}
		val := modPow(a%Q, exp, Q)
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, val)
	}
	out.WriteByte('\n')
}
