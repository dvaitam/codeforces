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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		ans := 0
		i := 0
		for i < n-1 {
			ans++
			j := i
			if a[j+1] > a[j] {
				for j+1 < n && a[j+1] > a[j] {
					j++
				}
			} else {
				for j+1 < n && a[j+1] < a[j] {
					j++
				}
			}
			i = j
		}
		ans++
		fmt.Fprintln(out, ans)
	}
}
