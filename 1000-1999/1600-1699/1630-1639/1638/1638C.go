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
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		maxVal := 0
		comps := 0
		for i := 0; i < n; i++ {
			if p[i] > maxVal {
				maxVal = p[i]
			}
			if maxVal == i+1 {
				comps++
			}
		}
		fmt.Fprintln(out, comps)
	}
}
