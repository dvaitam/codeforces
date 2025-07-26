package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var a1, x, y, m, k int
	if _, err := fmt.Fscan(reader, &n, &a1, &x, &y, &m, &k); err != nil {
		return
	}

	a := int64(a1)
	xm := int64(x)
	ym := int64(y)
	mm := int64(m)

	b := make([]int64, k+1)
	var ans uint64

	for i := 1; i <= n; i++ {
		ai := a
		for t := k; t >= 1; t-- {
			val := (b[t] + b[t-1]) % mod
			if t == 1 {
				val = (val + ai) % mod
			}
			b[t] = val
		}
		b[0] = (b[0] + ai) % mod

		ans ^= uint64(b[k]) * uint64(i)

		a = (a*xm + ym) % mm
	}

	fmt.Fprintln(writer, ans)
}
