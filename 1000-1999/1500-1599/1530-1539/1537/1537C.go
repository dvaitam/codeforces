package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		h := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &h[i])
		}
		sort.Ints(h)
		if n == 2 {
			fmt.Fprintf(writer, "%d %d\n", h[0], h[1])
			continue
		}
		idx := 0
		minDiff := h[1] - h[0]
		for i := 1; i < n-1; i++ {
			if d := h[i+1] - h[i]; d < minDiff {
				minDiff = d
				idx = i
			}
		}
		res := make([]int, 0, n)
		res = append(res, h[idx])
		for i := idx + 2; i < n; i++ {
			res = append(res, h[i])
		}
		for i := 0; i < idx; i++ {
			res = append(res, h[i])
		}
		res = append(res, h[idx+1])

		for i, v := range res {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
