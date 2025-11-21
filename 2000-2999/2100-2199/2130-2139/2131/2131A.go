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
		sort.Ints(sort.IntSlice(a))
		sort.Ints(sort.IntSlice(b))
		for i := n / 2; i < n; i++ {
			a[i]++
		}
		ans := 0
		for i := 0; i < n; i++ {
			if a[i] >= b[i] {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}

