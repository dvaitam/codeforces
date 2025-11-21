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
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		var sumAbs int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			if x < 0 {
				sumAbs -= x
			} else {
				sumAbs += x
			}
		}
		fmt.Fprintln(out, sumAbs)
	}
}
