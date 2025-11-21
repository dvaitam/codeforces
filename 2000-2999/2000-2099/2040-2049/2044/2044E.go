package main

import (
	"bufio"
	"fmt"
	"os"
)

func ceilDiv(a, b int64) int64 {
	return (a + b - 1) / b
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
		var k, l1, r1, l2, r2 int64
		fmt.Fscan(in, &k, &l1, &r1, &l2, &r2)

		var ans int64
		p := int64(1)
		for {
			left := l1
			if x := ceilDiv(l2, p); x > left {
				left = x
			}
			right := r1
			if y := r2 / p; y < right {
				right = y
			}
			if left <= right {
				ans += right - left + 1
			}
			if p > r2/k {
				break
			}
			p *= k
		}

		fmt.Fprintln(out, ans)
	}
}
