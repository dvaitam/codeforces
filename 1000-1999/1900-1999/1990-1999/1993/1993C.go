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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		L := 2 * k
		diff := make([]int, L+1)
		maxA := 0
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(in, &a)
			if a > maxA {
				maxA = a
			}
			b := a % L
			end := b + k
			if end <= L {
				diff[b]++
				diff[end]--
			} else {
				diff[b]++
				diff[L]--
				endWrap := end - L
				diff[0]++
				diff[endWrap]--
			}
		}

		w := -1
		cur := 0
		for i := 0; i < L; i++ {
			cur += diff[i]
			if cur == n {
				w = i
				break
			}
		}

		if w == -1 {
			fmt.Fprintln(out, -1)
			continue
		}

		rem := maxA % L
		ans := maxA - rem + w
		if ans < maxA {
			ans += L
		}
		fmt.Fprintln(out, ans)
	}
}
