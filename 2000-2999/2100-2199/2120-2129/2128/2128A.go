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
		var c int64
		fmt.Fscan(in, &n, &c)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		limits := make([]int, 0, n)
		cost := int64(0)
		for _, val := range a {
			if val > c {
				cost++
				continue
			}
			lim := 0
			cur := val
			for cur <= c {
				lim++
				if cur > c/2 {
					break
				}
				cur *= 2
			}
			limits = append(limits, lim)
		}
		sort.Ints(limits)
		time := 1
		free := 0
		for _, lim := range limits {
			if lim >= time {
				free++
				time++
			}
		}
		ans := cost + int64(len(limits)-free)
		fmt.Fprintln(out, ans)
	}
}
