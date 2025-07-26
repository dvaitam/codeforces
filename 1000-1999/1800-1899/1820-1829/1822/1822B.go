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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		p1 := a[0] * a[1]
		p2 := a[n-1] * a[n-2]
		if p1 > p2 {
			fmt.Fprintln(out, p1)
		} else {
			fmt.Fprintln(out, p2)
		}
	}
}
