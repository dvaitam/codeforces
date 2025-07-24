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
		for i := range a {
			fmt.Fscan(in, &a[i])
		}
		win := false
		mn := a[0]
		for i := 1; i < n; i++ {
			if a[i] < mn {
				win = true
				break
			}
		}
		if win {
			fmt.Fprintln(out, "Alice")
		} else {
			fmt.Fprintln(out, "Bob")
		}
	}
}
