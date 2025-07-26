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
		fmt.Fscan(in, &n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, 2*i)
		}
		fmt.Fprintln(out)
	}
}
