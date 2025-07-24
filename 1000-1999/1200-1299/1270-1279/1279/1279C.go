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
		var n, m int
		fmt.Fscan(in, &n, &m)

		a := make([]int, n)
		pos := make(map[int]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			pos[a[i]] = i
		}

		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}

		var res int64
		mx := -1
		for i, x := range b {
			p := pos[x]
			if p <= mx {
				res += 1
			} else {
				res += int64(2*(p-i) + 1)
				mx = p
			}
		}
		fmt.Fprintln(out, res)
	}
}
