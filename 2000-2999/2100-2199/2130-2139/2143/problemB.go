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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		var total int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			total += a[i]
		}
		b := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &b[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		sort.Ints(b)
		pos := 0
		for _, x := range b {
			if pos+x > n {
				break
			}
			pos += x
			total -= a[pos-1]
		}
		fmt.Fprintln(out, total)
	}
}
