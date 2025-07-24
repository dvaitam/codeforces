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
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
		}
		k := 10 - n
		ans := 3 * k * (k - 1)
		fmt.Fprintln(out, ans)
	}
}
