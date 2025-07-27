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
		var k int
		fmt.Fscan(in, &n, &k)
		h := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}
		pos := -1
		for step := 0; step < k; step++ {
			idx := -1
			for i := 0; i < n-1; i++ {
				if h[i] < h[i+1] {
					idx = i
					break
				}
			}
			if idx == -1 {
				pos = -1
				break
			}
			h[idx]++
			pos = idx + 1
		}
		fmt.Fprintln(out, pos)
	}
}
