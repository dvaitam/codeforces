package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(a []int, l, r, depth int, res []int) {
	if l > r {
		return
	}
	maxIdx := l
	for i := l + 1; i <= r; i++ {
		if a[i] > a[maxIdx] {
			maxIdx = i
		}
	}
	res[maxIdx] = depth
	solve(a, l, maxIdx-1, depth+1, res)
	solve(a, maxIdx+1, r, depth+1, res)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		res := make([]int, n)
		solve(a, 0, n-1, 0, res)
		for i, v := range res {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		writer.WriteByte('\n')
	}
}
