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
		var a, b int
		fmt.Fscan(in, &a, &b)
		if a == 0 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, a+2*b+1)
		}
	}
}
