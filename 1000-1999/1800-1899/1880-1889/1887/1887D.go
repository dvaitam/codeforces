package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		r--
		good := false
		for i := l; i < r && !good; i++ {
			maxLeft := a[l]
			for j := l; j <= i; j++ {
				if a[j] > maxLeft {
					maxLeft = a[j]
				}
			}
			minRight := a[i+1]
			for j := i + 1; j <= r; j++ {
				if a[j] < minRight {
					minRight = a[j]
				}
			}
			if maxLeft < minRight {
				good = true
			}
		}
		if good {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
