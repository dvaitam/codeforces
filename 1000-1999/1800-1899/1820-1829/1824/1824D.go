package main

import (
	"bufio"
	"fmt"
	"os"
)

func g(a []int, i, j int) int {
	if i > j {
		return 0
	}
	required := make(map[int]struct{})
	for p := i; p <= j; p++ {
		required[a[p-1]] = struct{}{}
	}
	for x := j; x >= 1; x-- {
		if _, ok := required[a[x-1]]; ok {
			delete(required, a[x-1])
			if len(required) == 0 {
				return x
			}
		}
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; q > 0; q-- {
		var l, r, x, y int
		fmt.Fscan(in, &l, &r, &x, &y)
		ans := 0
		for i := l; i <= r; i++ {
			for j := x; j <= y; j++ {
				if i <= j {
					ans += g(a, i, j)
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
