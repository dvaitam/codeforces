package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func countTwos(x int64) int {
	c := 0
	for x%2 == 0 {
		c++
		x /= 2
	}
	return c
}

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
		fmt.Fscan(in, &n)
		cnt := 0
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			cnt += countTwos(v)
		}
		if cnt >= n {
			fmt.Fprintln(out, 0)
			continue
		}
		need := n - cnt
		powers := make([]int, n)
		for i := 1; i <= n; i++ {
			powers[i-1] = countTwos(int64(i))
		}
		sort.Slice(powers, func(i, j int) bool { return powers[i] > powers[j] })
		ops := 0
		for _, p := range powers {
			if need <= 0 {
				break
			}
			if p > 0 {
				need -= p
				ops++
			}
		}
		if need > 0 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ops)
		}
	}
}
