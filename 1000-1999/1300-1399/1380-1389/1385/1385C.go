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
		i := n - 1
		for i > 0 && a[i-1] >= a[i] {
			i--
		}
		for i > 0 && a[i-1] <= a[i] {
			i--
		}
		fmt.Fprintln(out, i)
	}
}
