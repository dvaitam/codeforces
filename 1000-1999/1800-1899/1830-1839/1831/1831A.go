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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}
		for i, v := range a {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, n+1-v)
		}
		fmt.Fprintln(out)
	}
}
