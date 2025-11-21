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
		ans := 0
		total := 2*n - 1
		for i := 0; i < total; i++ {
			var x int
			if _, err := fmt.Fscan(in, &x); err != nil {
				return
			}
			ans ^= x
		}
		fmt.Fprintln(out, ans)
	}
}
