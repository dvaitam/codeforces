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
		var s string
		fmt.Fscan(in, &s)

		improvements := make([]int64, n)
		var total int64
		for i := 0; i < n; i++ {
			if s[i] == 'L' {
				total += int64(i)
				improvements[i] = int64(n - 1 - 2*i)
			} else {
				total += int64(n - 1 - i)
				improvements[i] = int64(2*i - n + 1)
			}
		}

		sort.Slice(improvements, func(i, j int) bool { return improvements[i] > improvements[j] })

		prefix := int64(0)
		for i := 0; i < n; i++ {
			if improvements[i] > 0 {
				prefix += improvements[i]
			}
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, total+prefix)
		}
		out.WriteByte('\n')
	}
}
