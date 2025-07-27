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
		prefix := 0
		for prefix < n && a[prefix] == 1 {
			prefix++
		}
		if prefix == n {
			if n%2 == 1 {
				fmt.Fprintln(out, "First")
			} else {
				fmt.Fprintln(out, "Second")
			}
		} else {
			if prefix%2 == 0 {
				fmt.Fprintln(out, "First")
			} else {
				fmt.Fprintln(out, "Second")
			}
		}
	}
}
