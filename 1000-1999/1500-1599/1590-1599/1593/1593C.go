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
		var n int64
		var k int
		fmt.Fscan(in, &n, &k)
		xs := make([]int64, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &xs[i])
		}
		sort.Slice(xs, func(i, j int) bool { return xs[i] > xs[j] })
		var dist int64
		saved := 0
		for _, x := range xs {
			need := n - x
			if dist+need < n {
				saved++
				dist += need
			} else {
				break
			}
		}
		fmt.Fprintln(out, saved)
	}
}
