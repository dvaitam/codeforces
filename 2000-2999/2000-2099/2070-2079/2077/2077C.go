package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

var pow2 []int
var inv4 int

func modPow(a, e int) int {
	res := 1
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		e >>= 1
	}
	return res
}

func prepare(maxN int) {
	pow2 = make([]int, maxN+1)
	pow2[0] = 1
	for i := 1; i <= maxN; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}
	inv4 = modPow(4, mod-2)
}

// scoreSum computes the required sum knowing only n and the count of ones.
// Any subsequence is determined by choosing subsets of ones and zeros,
// so the final value depends solely on these counts, which leads to the
// closed-form formula derived in the analysis.
func scoreSum(n, ones int) int {
	zeros := n - ones
	diff := ones - zeros
	diffMod := diff % mod
	if diffMod < 0 {
		diffMod += mod
	}
	term := (int(int64(diffMod)*int64(diffMod)%mod) + n) % mod
	pow2n := pow2[n]
	A := int(int64(pow2n) * int64(inv4) % mod * int64(term) % mod)
	B := pow2[n-1]
	res := A - B
	if res < 0 {
		res += mod
	}
	res = int(int64(res) * int64(inv4) % mod)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	maxN := 200000 + 5
	prepare(maxN)

	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		var s string
		fmt.Fscan(in, &s)
		bytes := []byte(s)
		ones := 0
		for _, ch := range bytes {
			if ch == '1' {
				ones++
			}
		}

		for ; q > 0; q-- {
			var idx int
			fmt.Fscan(in, &idx)
			idx--
			if bytes[idx] == '1' {
				bytes[idx] = '0'
				ones--
			} else {
				bytes[idx] = '1'
				ones++
			}
			fmt.Fprintln(out, scoreSum(n, ones))
		}
	}
}
