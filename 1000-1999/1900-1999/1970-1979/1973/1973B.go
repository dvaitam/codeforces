package main

import (
	"bufio"
	"fmt"
	"os"
)

func loneliness(a []int) int {
	n := len(a)
	maxGap := 0
	for bit := 0; bit < 20; bit++ {
		last := 0
		appear := false
		for i := 0; i < n; i++ {
			if (a[i]>>bit)&1 == 1 {
				appear = true
				gap := i + 1 - last - 1
				if gap > maxGap {
					maxGap = gap
				}
				last = i + 1
			}
		}
		if appear {
			gap := n + 1 - last - 1
			if gap > maxGap {
				maxGap = gap
			}
		}
	}
	return maxGap + 1
}

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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, loneliness(a))
	}
}
