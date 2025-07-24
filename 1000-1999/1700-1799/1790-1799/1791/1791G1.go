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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var c int64
		fmt.Fscan(in, &n, &c)
		costs := make([]int64, n)
		for i := 0; i < n; i++ {
			var a int64
			fmt.Fscan(in, &a)
			costs[i] = a + int64(i+1)
		}
		sort.Slice(costs, func(i, j int) bool { return costs[i] < costs[j] })
		count := 0
		for _, v := range costs {
			if c >= v {
				c -= v
				count++
			} else {
				break
			}
		}
		fmt.Fprintln(out, count)
	}
}
