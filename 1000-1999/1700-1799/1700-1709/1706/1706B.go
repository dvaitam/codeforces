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
		odd := make([]int, n+1)
		even := make([]int, n+1)
		for i := 1; i <= n; i++ {
			var c int
			fmt.Fscan(in, &c)
			if i%2 == 1 {
				odd[c]++
			} else {
				even[c]++
			}
		}
		for r := 1; r <= n; r++ {
			if odd[r] > even[r] {
				fmt.Fprint(out, odd[r])
			} else {
				fmt.Fprint(out, even[r])
			}
			if r < n {
				fmt.Fprint(out, " ")
			} else {
				fmt.Fprintln(out)
			}
		}
	}
}
