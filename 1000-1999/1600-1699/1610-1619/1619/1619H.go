package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &p[i])
	}

	B := 1
	for B*B < n {
		B++
	}
	block := make([]int, n+1)
	for i := 1; i <= n; i++ {
		block[i] = (i - 1) / B
	}
	jump := make([]int, n+1)
	cnt := make([]int, n+1)

	rebuild := func(b int) {
		l := b*B + 1
		r := (b + 1) * B
		if r > n {
			r = n
		}
		for i := r; i >= l; i-- {
			nxt := p[i]
			if nxt >= l && nxt <= r {
				jump[i] = jump[nxt]
				cnt[i] = cnt[nxt] + 1
			} else {
				jump[i] = nxt
				cnt[i] = 1
			}
		}
	}

	blocks := (n + B - 1) / B
	for b := 0; b < blocks; b++ {
		rebuild(b)
	}

	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(in, &t, &x, &y)
		if t == 1 {
			p[x], p[y] = p[y], p[x]
			rebuild(block[x])
			if block[x] != block[y] {
				rebuild(block[y])
			}
		} else {
			i := x
			k := y
			for k > 0 {
				if cnt[i] <= k {
					k -= cnt[i]
					i = jump[i]
				} else {
					i = p[i]
					k--
				}
			}
			fmt.Fprintln(out, i)
		}
	}
}
