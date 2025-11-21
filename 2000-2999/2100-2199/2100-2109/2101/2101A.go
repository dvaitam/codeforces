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
		cur := 0
		for i := 0; i < n; i++ {
			row := make([]int, n)
			if i%2 == 0 {
				for j := 0; j < n; j++ {
					row[j] = cur
					cur++
				}
			} else {
				for j := n - 1; j >= 0; j-- {
					row[j] = cur
					cur++
				}
			}
			for j := 0; j < n; j++ {
				if j > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, row[j])
			}
			fmt.Fprintln(out)
		}
	}
}
