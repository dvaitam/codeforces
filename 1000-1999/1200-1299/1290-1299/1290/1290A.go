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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if k > m-1 {
			k = m - 1
		}
		uncontrolled := m - 1 - k
		finalLen := n - m + 1
		ans := 0
		for takeFront := 0; takeFront <= k; takeFront++ {
			best := int(^uint(0) >> 1)
			for skipFront := 0; skipFront <= uncontrolled; skipFront++ {
				l := takeFront + skipFront
				r := l + finalLen - 1
				cand := a[l]
				if a[r] > cand {
					cand = a[r]
				}
				if cand < best {
					best = cand
				}
			}
			if best > ans {
				ans = best
			}
		}
		fmt.Fprintln(out, ans)
	}
}
