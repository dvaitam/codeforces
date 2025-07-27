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
		a := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Fscan(in, &a[i])
		}
		seen := make([]bool, n+1)
		res := make([]int, 0, n)
		for _, v := range a {
			if !seen[v] {
				res = append(res, v)
				seen[v] = true
			}
		}
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
