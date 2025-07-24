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

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Ints(a)
		teams := 1
		for i := 1; i < n; i++ {
			if a[i]-a[i-1] == 1 {
				teams = 2
				break
			}
		}
		fmt.Fprintln(out, teams)
	}
}
