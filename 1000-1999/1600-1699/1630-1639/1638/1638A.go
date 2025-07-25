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
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		for i := 0; i < n; i++ {
			if p[i] != i+1 {
				j := i
				for ; j < n; j++ {
					if p[j] == i+1 {
						break
					}
				}
				for l, r := i, j; l < r; l, r = l+1, r-1 {
					p[l], p[r] = p[r], p[l]
				}
				break
			}
		}
		for i, v := range p {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
