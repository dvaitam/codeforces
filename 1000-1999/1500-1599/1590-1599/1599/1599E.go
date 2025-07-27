package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7

// fast doubling Fibonacci implementation
func fibPair(n int64) (int64, int64) {
	if n == 0 {
		return 0, 1
	}
	a, b := fibPair(n >> 1)
	c := a * ((b*2 - a + MOD) % MOD) % MOD
	d := (a*a%MOD + b*b%MOD) % MOD
	if n&1 == 0 {
		return c, d
	}
	return d, (c + d) % MOD
}

func fib(n int64) int64 {
	x, _ := fibPair(n)
	return x % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	a1 := make([]int64, n)
	a2 := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a1[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a2[i])
	}

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 4 {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l--
			r--
			var sum int64
			for i := l; i <= r; i++ {
				sum += fib(a1[i] + a2[i])
				if sum >= MOD {
					sum %= MOD
				}
			}
			fmt.Fprintln(writer, sum%MOD)
		} else {
			var k, l, r int
			var x int64
			fmt.Fscan(reader, &k, &l, &r, &x)
			l--
			r--
			arr := &a1
			if k == 2 {
				arr = &a2
			}
			switch t {
			case 1:
				for i := l; i <= r; i++ {
					if (*arr)[i] > x {
						(*arr)[i] = x
					}
				}
			case 2:
				for i := l; i <= r; i++ {
					if (*arr)[i] < x {
						(*arr)[i] = x
					}
				}
			case 3:
				for i := l; i <= r; i++ {
					(*arr)[i] += x
				}
			}
		}
	}
}
