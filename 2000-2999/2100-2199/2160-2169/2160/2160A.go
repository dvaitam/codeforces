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

		freq := make([]int, 102)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			freq[x]++
		}

		mex := 0
		for ; mex < len(freq); mex++ {
			if freq[mex] == 0 {
				break
			}
		}
		fmt.Fprintln(out, mex)
	}
}
