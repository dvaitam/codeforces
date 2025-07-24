package main

import (
	"bufio"
	"fmt"
	"os"
)

func modPow(a, e, mod int64) int64 {
	res := int64(1)
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	var t int64
	var w int64
	if _, err := fmt.Fscan(in, &n, &m, &t, &w); err != nil {
		return
	}

	values := make(map[int]int64)
	var knownXor int64
	for i := 0; i < m; i++ {
		var d int
		var e int64
		fmt.Fscan(in, &d, &e)
		values[d] = e
		knownXor ^= e
	}
	lost := int64(n - len(values))

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var f int
		var g int64
		var mod int64
		fmt.Fscan(in, &f, &g, &mod)
		if old, ok := values[f]; ok {
			knownXor ^= old
			delete(values, f)
			lost++
		}
		if g != -1 {
			values[f] = g
			knownXor ^= g
			lost--
		}

		if lost == 0 {
			if knownXor%(1<<w) == 0 {
				fmt.Fprintln(out, 1%mod)
			} else {
				fmt.Fprintln(out, 0)
			}
			continue
		}
		exp := w * (lost - 1)
		ans := modPow(2, exp, mod)
		fmt.Fprintln(out, ans)
	}
}
