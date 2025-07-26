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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if k == 1 {
			fmt.Fprintln(out, (n-1)/2)
		} else {
			cnt := 0
			for i := 1; i < n-1; i++ {
				if a[i] > a[i-1]+a[i+1] {
					cnt++
				}
			}
			fmt.Fprintln(out, cnt)
		}
	}
}
