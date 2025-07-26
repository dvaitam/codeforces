package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

type pair struct {
	s int
	m int
}

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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		p := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			var x, y int64
			fmt.Fscan(reader, &x, &y)
			p[i] = x * modInv(y) % MOD
		}
		h := make([]int64, n+1)
		for i := 0; i <= n; i++ {
			fmt.Fscan(reader, &h[i])
		}

		dp := map[pair]int64{{0, 0}: 1}
		results := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			ndp := make(map[pair]int64)
			pi := p[i]
			qi := (1 + MOD - pi) % MOD
			for st, val := range dp {
				// step +1
				s1 := st.s + 1
				m1 := st.m
				if s1 > m1 {
					m1 = s1
				}
				ndp[pair{s1, m1}] = (ndp[pair{s1, m1}] + val*pi) % MOD
				// step -1
				s2 := st.s - 1
				m2 := st.m
				ndp[pair{s2, m2}] = (ndp[pair{s2, m2}] + val*qi) % MOD
			}
			dp = ndp
			var exp int64
			for st, v := range dp {
				exp = (exp + v*h[st.m]) % MOD
			}
			results[i] = exp
		}
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, results[i])
		}
		fmt.Fprintln(writer)
	}
}
