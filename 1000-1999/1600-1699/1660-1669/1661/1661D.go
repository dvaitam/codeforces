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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	diff1 := make([]int64, n+2)
	diff2 := make([]int64, n+2)
	var cur1, cur2 int64
	var ans int64

	for i := n; i >= 1; i-- {
		cur1 += diff1[i]
		cur2 += diff2[i]
		val := cur1*int64(i) + cur2
		if val < b[i-1] {
			delta := b[i-1] - val
			if i >= k {
				x := (delta + int64(k) - 1) / int64(k)
				ans += x
				l := i - k + 1
				cur1 += x
				cur2 += x * (1 - int64(l))
				diff1[l-1] -= x
				diff2[l-1] -= x * (1 - int64(l))
			} else {
				x := (delta + int64(i) - 1) / int64(i)
				ans += x
				cur1 += x
			}
		}
	}

	fmt.Fprintln(out, ans)
}
