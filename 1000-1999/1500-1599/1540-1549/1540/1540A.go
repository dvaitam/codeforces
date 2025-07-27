package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		d := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &d[i])
		}
		sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
		var sum int64
		for i := 0; i < n; i++ {
			coef := int64(n - 2*i - 1)
			sum += d[i] * coef
		}
		result := d[n-1] + sum
		fmt.Fprintln(out, result)
	}
}
