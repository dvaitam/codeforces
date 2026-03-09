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

	var n int
	var A, B int64
	fmt.Fscan(in, &n, &A, &B)

	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 0; i < n; i++ {
		var ch string
		fmt.Fscan(in, &a[i], &ch)
		if ch == "W" {
			b[i] = 1
		} else {
			b[i] = 0
		}
	}
	fmt.Fscan(in, &a[n])
	b[n] = 2 // sentinel: neither W nor P

	s := int64(n) * B
	var t, k int64
	for i := n - 1; i >= 0; i-- {
		t += (a[i+1] - a[i]) * A
		if b[i] == b[i+1] {
			v := (k - a[i+1]) * A
			if v > B {
				v = B
			}
			t += v
		} else {
			k = a[i+1]
		}
		v := t + int64(i)*B
		if v < s {
			s = v
		}
	}
	fmt.Fprintln(out, s)
}
