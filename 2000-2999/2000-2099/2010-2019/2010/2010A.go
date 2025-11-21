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
		sign := 1
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			sum += sign * x
			sign *= -1
		}
		fmt.Fprintln(out, sum)
	}
}
