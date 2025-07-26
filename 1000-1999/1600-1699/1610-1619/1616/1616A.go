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
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x < 0 {
				x = -x
			}
			freq[x]++
		}
		res := 0
		for v, c := range freq {
			if v == 0 {
				res += 1
			} else {
				if c > 2 {
					res += 2
				} else {
					res += c
				}
			}
		}
		fmt.Fprintln(out, res)
	}
}
