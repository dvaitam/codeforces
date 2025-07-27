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
		}
		ans := []int{p[0]}
		for i := 1; i < n-1; i++ {
			if (p[i]-p[i-1])*(p[i+1]-p[i]) < 0 {
				ans = append(ans, p[i])
			}
		}
		ans = append(ans, p[n-1])
		fmt.Fprintln(out, len(ans))
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
