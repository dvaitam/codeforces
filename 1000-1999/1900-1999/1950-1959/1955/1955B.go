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
		var c, d int64
		fmt.Fscan(in, &n, &c, &d)
		total := n * n
		b := make([]int64, total)
		for i := 0; i < total; i++ {
			fmt.Fscan(in, &b[i])
		}
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

		base := b[0]
		expected := make([]int64, total)
		idx := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				expected[idx] = base + int64(i)*c + int64(j)*d
				idx++
			}
		}
		sort.Slice(expected, func(i, j int) bool { return expected[i] < expected[j] })
		ok := true
		for i := 0; i < total && ok; i++ {
			if expected[i] != b[i] {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
