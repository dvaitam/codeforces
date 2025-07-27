package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement a correct solution
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, -1)
		}
		fmt.Fprintln(out)
	}
}
