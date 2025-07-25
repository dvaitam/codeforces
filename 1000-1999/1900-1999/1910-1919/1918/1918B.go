package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		idx := make([]int, n)
		for i := 0; i < n; i++ {
			idx[i] = i
		}
		sort.Slice(idx, func(i, j int) bool {
			return a[idx[i]] < a[idx[j]]
		})
		// output a' and b'
		for i, id := range idx {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, a[id])
		}
		out.WriteByte('\n')
		for i, id := range idx {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, b[id])
		}
		out.WriteByte('\n')
	}
}
