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
		sum := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x < 0 {
				sum -= x
			} else {
				sum += x
			}
		}
		fmt.Fprintln(out, sum)
	}
}
