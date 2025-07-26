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
		var x int
		fmt.Fscan(in, &x)
		var y int
		if x == 1 {
			y = 3
		} else if x&1 == 1 {
			y = 1
		} else if x&(x-1) == 0 {
			y = x + 1
		} else {
			y = x & -x
		}
		fmt.Fprintln(out, y)
	}
}
