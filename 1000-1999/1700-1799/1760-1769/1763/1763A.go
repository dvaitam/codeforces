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
		var orVal, andVal int
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if i == 0 {
				orVal = x
				andVal = x
			} else {
				orVal |= x
				andVal &= x
			}
		}
		fmt.Fprintln(out, orVal-andVal)
	}
}
