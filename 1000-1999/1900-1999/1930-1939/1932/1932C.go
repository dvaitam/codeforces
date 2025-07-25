package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var s string
		fmt.Fscan(reader, &s)
		removed := make([]int, n)
		l, r := 0, n-1
		for i := 0; i < n; i++ {
			if s[i] == 'L' {
				removed[i] = a[l]
				l++
			} else {
				removed[i] = a[r]
				r--
			}
		}
		ans := make([]int, n)
		prod := 1 % m
		for i := n - 1; i >= 0; i-- {
			prod = (prod * removed[i]) % m
			ans[i] = prod
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[i])
		}
		fmt.Fprintln(writer)
	}
}
