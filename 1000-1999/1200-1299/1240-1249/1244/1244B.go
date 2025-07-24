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
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		ans := n
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				left := i + 1
				right := n - i
				if left > right {
					if 2*left > ans {
						ans = 2 * left
					}
				} else {
					if 2*right > ans {
						ans = 2 * right
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
