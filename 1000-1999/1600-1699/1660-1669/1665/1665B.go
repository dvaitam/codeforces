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
		freq := make(map[int]int, n)
		mx := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			freq[x]++
			if freq[x] > mx {
				mx = freq[x]
			}
		}
		ops := 0
		for mx < n {
			ops++
			diff := n - mx
			add := mx
			if diff < add {
				add = diff
			}
			ops += add
			mx += add
		}
		fmt.Fprintln(out, ops)
	}
}
