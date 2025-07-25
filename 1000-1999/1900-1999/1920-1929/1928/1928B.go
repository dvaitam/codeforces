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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Ints(arr)
		cnt := make(map[int]int)
		l := 0
		best := 0
		for r, v := range arr {
			cnt[v]++
			for arr[r]-arr[l] > n-1 {
				val := arr[l]
				cnt[val]--
				if cnt[val] == 0 {
					delete(cnt, val)
				}
				l++
			}
			if len(cnt) > best {
				best = len(cnt)
			}
		}
		fmt.Fprintln(out, best)
	}
}
