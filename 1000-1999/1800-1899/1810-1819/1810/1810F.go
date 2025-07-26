package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func computeDepths(n, m int) []int {
	q := []int{0}
	for len(q) < n {
		d := q[0]
		q = q[1:]
		for i := 0; i < m; i++ {
			q = append(q, d+1)
		}
	}
	return q[:n]
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
		var n, m, q int
		fmt.Fscan(reader, &n, &m, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		depth := computeDepths(n, m)
		for ; q > 0; q-- {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			x--
			a[x] = y
			b := make([]int, n)
			copy(b, a)
			sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
			res := 0
			for i := 0; i < n; i++ {
				v := b[i] + depth[i]
				if v > res {
					res = v
				}
			}
			fmt.Fprint(writer, res)
			if q > 0 {
				fmt.Fprint(writer, " ")
			}
		}
		fmt.Fprintln(writer)
	}
}
