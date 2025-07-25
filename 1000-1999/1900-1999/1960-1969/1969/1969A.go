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
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			p[i]--
		}
		ans := 3
		for i := 0; i < n; i++ {
			if p[p[i]] == i {
				ans = 2
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
