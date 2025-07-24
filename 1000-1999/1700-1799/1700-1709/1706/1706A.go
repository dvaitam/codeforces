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
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		res := make([]byte, m)
		for i := 0; i < m; i++ {
			res[i] = 'B'
		}
		for _, val := range a {
			left := val
			right := m + 1 - val
			if left > right {
				left, right = right, left
			}
			if res[left-1] == 'B' {
				res[left-1] = 'A'
			} else {
				res[right-1] = 'A'
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
