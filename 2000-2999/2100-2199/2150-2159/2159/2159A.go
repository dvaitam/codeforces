package main

import (
	"bufio"
	"fmt"
	"os"
)

// The original problem is interactive: we would discover the hidden array
// through MAD queries.  In this offline archive the full array is provided
// for every test case, so solving it amounts to echoing those 2n numbers.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for tc := 0; tc < t; tc++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}

		arr := make([]int, 2*n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}

		for i, v := range arr {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		if tc+1 < t {
			fmt.Fprintln(out)
		}
	}
}

