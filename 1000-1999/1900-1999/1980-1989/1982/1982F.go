package main

import (
	"bufio"
	"fmt"
	"os"
)

func findSegment(a []int) (int, int) {
	n := len(a)
	l := 0
	for l+1 < n && a[l] <= a[l+1] {
		l++
	}
	if l == n-1 {
		return -1, -1
	}
	r := n - 1
	for r > 0 && a[r-1] <= a[r] {
		r--
	}
	minV, maxV := a[l], a[l]
	for i := l; i <= r; i++ {
		if a[i] < minV {
			minV = a[i]
		}
		if a[i] > maxV {
			maxV = a[i]
		}
	}
	for l > 0 && a[l-1] > minV {
		l--
	}
	for r < n-1 && a[r+1] < maxV {
		r++
	}
	return l + 1, r + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		var q int
		fmt.Fscan(in, &q)
		l, r := findSegment(a)
		fmt.Fprintln(out, l, r)
		for ; q > 0; q-- {
			var pos, val int
			fmt.Fscan(in, &pos, &val)
			a[pos-1] = val
			l, r = findSegment(a)
			fmt.Fprintln(out, l, r)
		}
	}
}
