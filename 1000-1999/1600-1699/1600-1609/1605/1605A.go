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
		var a1, a2, a3 int
		fmt.Fscan(in, &a1, &a2, &a3)
		sum := a1 + a2 + a3
		if sum%3 == 0 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, 1)
		}
	}
}
