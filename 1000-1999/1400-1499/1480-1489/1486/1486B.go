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
		xs := make([]int64, n)
		ys := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &xs[i], &ys[i])
		}
		sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
		sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
		if n%2 == 1 {
			fmt.Fprintln(writer, 1)
		} else {
			xRange := xs[n/2] - xs[n/2-1] + 1
			yRange := ys[n/2] - ys[n/2-1] + 1
			fmt.Fprintln(writer, xRange*yRange)
		}
	}
}
