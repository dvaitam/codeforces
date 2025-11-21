package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const mod int64 = 1000000007

func fastPow2(exp int) int64 {
	if exp == 0 {
		return 1
	}
	base := int64(2)
	res := int64(1)
	e := exp
	for e > 0 {
		if e&1 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		e >>= 1
	}
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

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		sumOdd := int64(0)
		pow2 := int64(1)
		for i := 1; i <= n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			cnt := 0
			for x%2 == 0 {
				x >>= 1
				cnt++
			}
			sumOdd = (sumOdd + x) % mod
			pow2 = (pow2 * fastPow2(cnt)) % mod
			ans := (sumOdd - x) % mod
			if ans < 0 {
				ans += mod
			}
			ans = (ans + (x%mod)*pow2) % mod
			if i > 1 {
				out.WriteByte(' ')
			}
			out.WriteString(strconv.FormatInt(ans, 10))
		}
		out.WriteByte('\n')
	}
}
