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
		freq := make([]int, n+1)
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(in, &a)
			if a <= n {
				freq[a]++
			}
		}
		count := make([]int, n+1)
		maxCatch := 0
		for d := 1; d <= n; d++ {
			if freq[d] == 0 {
				continue
			}
			for m := d; m <= n; m += d {
				count[m] += freq[d]
				if count[m] > maxCatch {
					maxCatch = count[m]
				}
			}
		}
		fmt.Fprintln(out, maxCatch)
	}
}
