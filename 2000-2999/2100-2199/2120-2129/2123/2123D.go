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
		var s string
		fmt.Fscan(in, &s)

		left, right := 0, n-1
		for left < n && s[left] == '0' {
			left++
		}
		for right >= 0 && s[right] == '0' {
			right--
		}

		if left > right {
			fmt.Fprintln(out, "Alice")
			continue
		}

		span := right - left + 1
		if span <= k {
			fmt.Fprintln(out, "Alice")
		} else if span > k+1 {
			fmt.Fprintln(out, "Bob")
		} else {
			fmt.Fprintln(out, "Alice")
		}
	}
}
