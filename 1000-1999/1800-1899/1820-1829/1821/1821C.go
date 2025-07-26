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
		var s string
		fmt.Fscan(in, &s)
		n := len(s)
		best := int(1<<31 - 1)
		for ch := byte('a'); ch <= 'z'; ch++ {
			prev := -1
			maxSeg := 0
			for i := 0; i < n; i++ {
				if s[i] == ch {
					seg := i - prev - 1
					if seg > maxSeg {
						maxSeg = seg
					}
					prev = i
				}
			}
			seg := n - prev - 1
			if seg > maxSeg {
				maxSeg = seg
			}
			ops := 0
			for v := maxSeg; v > 0; v >>= 1 {
				ops++
			}
			if ops < best {
				best = ops
			}
		}
		fmt.Fprintln(out, best)
	}
}
