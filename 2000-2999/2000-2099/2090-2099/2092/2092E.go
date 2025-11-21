package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000007
const MAXK = 200000

var pow2 []int64

func initPow2() {
	pow2 = make([]int64, MAXK+5)
	pow2[0] = 1
	for i := 1; i < len(pow2); i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}
}

func main() {
	initPow2()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int64
		var k int
		fmt.Fscan(in, &n, &m, &k)

		parity := 0
		for i := 0; i < k; i++ {
			var x, y int64
			var c int
			fmt.Fscan(in, &x, &y, &c)
			cellParity := (x + y) & 1
			if int(cellParity) != c {
				parity ^= 1
			}
		}

		if parity == 1 {
			fmt.Fprintln(out, 0)
		} else {
			exponent := int64(n*m) - int64(k)
			if exponent < 0 {
				exponent = 0
			}
			if exponent > int64(MAXK) {
				fmt.Fprintln(out, modPow(2, exponent))
			} else {
				fmt.Fprintln(out, pow2[exponent])
			}
		}
	}
}

func modPow(base, exp int64) int64 {
	result := int64(1)
	b := base % MOD
	e := exp
	for e > 0 {
		if e&1 == 1 {
			result = (result * b) % MOD
		}
		b = (b * b) % MOD
		e >>= 1
	}
	return result
}
