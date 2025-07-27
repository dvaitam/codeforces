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
		s := make([]byte, n)
		for i := 0; i < n; i++ {
			s[i] = 'a'
		}
		for i := n - 2; i >= 0; i-- {
			cnt := n - i - 1
			if k > cnt {
				k -= cnt
			} else {
				s[i] = 'b'
				s[n-k] = 'b'
				break
			}
		}
		fmt.Fprintln(out, string(s))
	}
}
