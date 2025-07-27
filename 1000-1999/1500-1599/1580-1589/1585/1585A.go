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
		height := 1
		alive := true
		for i := 0; i < n && alive; i++ {
			if a[i] == 1 {
				if i > 0 && a[i-1] == 1 {
					height += 5
				} else {
					height += 1
				}
			} else {
				if i > 0 && a[i-1] == 0 {
					alive = false
				}
			}
		}
		if alive {
			fmt.Fprintln(out, height)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
