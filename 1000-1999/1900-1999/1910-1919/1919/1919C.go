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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		const inf = int(1e9 + 7)
		last := [2]int{inf, inf}
		penalty := 0
		for _, x := range a {
			best := -1
			bestVal := inf + 1
			for i := 0; i < 2; i++ {
				if last[i] >= x && last[i] <= bestVal {
					best = i
					bestVal = last[i]
				}
			}
			if best != -1 {
				last[best] = x
			} else {
				if last[0] <= last[1] {
					best = 0
				} else {
					best = 1
				}
				penalty++
				last[best] = x
			}
		}
		fmt.Fprintln(writer, penalty)
	}
}
