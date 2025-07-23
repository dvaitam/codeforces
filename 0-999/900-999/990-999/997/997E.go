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

	var n int
	fmt.Fscan(in, &n)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}
	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		r--
		var count int64
		for i := l; i <= r; i++ {
			minVal := p[i]
			maxVal := p[i]
			for j := i; j <= r; j++ {
				if p[j] < minVal {
					minVal = p[j]
				}
				if p[j] > maxVal {
					maxVal = p[j]
				}
				if maxVal-minVal == j-i {
					count++
				}
			}
		}
		fmt.Fprintln(out, count)
	}
}
