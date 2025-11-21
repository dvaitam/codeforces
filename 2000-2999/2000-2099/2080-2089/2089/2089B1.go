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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		type pair struct {
			a int64
			b int64
		}
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{a: a[i], b: b[i]}
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].b > arr[j].b
		})
		ans := int64(0)
		current := int64(0)
		for _, p := range arr {
			current += p.a
			if current > 0 {
				rounds := (current + p.b - 1) / p.b
				if rounds > ans {
					ans = rounds
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
