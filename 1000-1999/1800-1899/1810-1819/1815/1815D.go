package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

type state struct{ cnt, sum int64 }

func solveM2(n uint64) int64 {
	if n == 0 {
		return 0
	}
	bits := make([]int, 0)
	tmp := n
	for tmp > 0 {
		bits = append(bits, int(tmp&1))
		tmp >>= 1
	}
	bits = append(bits, 0) // sentinel
	L := len(bits)
	pow2 := make([]int64, L)
	pow2[0] = 1
	for i := 1; i < L; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	var dp [2][2]state
	dp[0][0] = state{1, 0}
	for i := 0; i < L; i++ {
		var ndp [2][2]state
		nbit := 0
		if i < L-1 {
			nbit = bits[i]
		}
		for b := 0; b < 2; b++ {
			for p := 0; p < 2; p++ {
				st := dp[b][p]
				if st.cnt == 0 {
					continue
				}
				temp := nbit - p - b
				var r, nb int
				if temp >= 0 {
					r = temp
					nb = 0
				} else {
					r = temp + 2
					nb = 1
				}
				if r == 0 {
					for y := 0; y <= 1; y++ {
						s := &ndp[nb][y]
						s.cnt = (s.cnt + st.cnt) % MOD
						s.sum = (s.sum + st.sum + st.cnt*((pow2[i]*int64(y))%MOD)) % MOD
					}
				} else {
					s := &ndp[nb][0]
					s.cnt = (s.cnt + st.cnt) % MOD
					s.sum = (s.sum + st.sum) % MOD
				}
			}
		}
		dp = ndp
	}
	cnt := dp[0][0].cnt % MOD
	sumY := dp[0][0].sum % MOD
	res := (cnt*(int64(n)%MOD) - 2*sumY) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n uint64
		var m int64
		fmt.Fscan(in, &n, &m)
		var ans int64
		if m == 1 {
			ans = int64(n % uint64(MOD))
		} else if m >= 3 {
			if n%2 == 0 {
				a := int64((n/2 + 1) % uint64(MOD))
				b := int64((n / 2) % uint64(MOD))
				ans = a * b % MOD
			} else {
				x := int64(((n + 1) / 2) % uint64(MOD))
				ans = x * x % MOD
			}
		} else { // m==2
			ans = solveM2(n)
		}
		fmt.Fprintln(out, ans)
	}
}
