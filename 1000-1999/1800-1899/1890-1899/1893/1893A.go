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
		var n, k int
		fmt.Fscan(in, &n, &k)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		s := 0
		steps := k
		if steps > n {
			steps = n
		}
		ok := true
		for i := 0; i < steps; i++ {
			idx := n - s - 1
			if idx < 0 {
				idx %= n
				idx += n
			}
			idx %= n
			x := b[idx]
			if x < 1 || x > n {
				ok = false
				break
			}
			s = (s + x) % n
		}
		if ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
