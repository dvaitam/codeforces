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
		times := make([]int, n)
		indices := make([]int, 0)
		for i := 0; i < n; i++ {
			if a[i] == i+1 {
				times[i] = 0
			} else {
				indices = append(indices, i)
			}
		}
		step := 0
		for len(indices) > 0 {
			step++
			if len(indices) > 0 {
				last := a[indices[len(indices)-1]]
				for i := len(indices) - 1; i > 0; i-- {
					a[indices[i]] = a[indices[i-1]]
				}
				a[indices[0]] = last
			}
			newIdx := make([]int, 0, len(indices))
			for _, idx := range indices {
				if a[idx] == idx+1 {
					times[idx] = step
				} else {
					newIdx = append(newIdx, idx)
				}
			}
			indices = newIdx
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, times[i])
		}
		fmt.Fprintln(writer)
	}
}
