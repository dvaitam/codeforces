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
	s := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}
	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(in, &typ)
		if typ == 1 {
			var i, t int
			fmt.Fscan(in, &i, &t)
			s[i-1] = t
		} else if typ == 2 {
			var l, r int
			fmt.Fscan(in, &l, &r)
			stack := make([]int, 0, r-l+1)
			valid := true
			for i := l - 1; i < r; i++ {
				v := s[i]
				if v > 0 {
					stack = append(stack, v)
				} else {
					t := -v
					if len(stack) == 0 || stack[len(stack)-1] != t {
						valid = false
						break
					}
					stack = stack[:len(stack)-1]
				}
			}
			if valid && len(stack) == 0 {
				fmt.Fprintln(out, "Yes")
			} else {
				fmt.Fprintln(out, "No")
			}
		}
	}
}
