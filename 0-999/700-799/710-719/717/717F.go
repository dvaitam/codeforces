package main

import (
	"bufio"
	"fmt"
	"os"
)

func canClearLeftRight(a []int) bool {
	if len(a) == 0 {
		return true
	}
	cur := a[0] - 1
	if cur < 0 {
		return false
	}
	for i := 1; i < len(a); i++ {
		cur = a[i] - 1 - cur
		if cur < 0 {
			return false
		}
	}
	return cur == a[len(a)-1]-1
}

func canClear(a []int) bool {
	if canClearLeftRight(a) {
		return true
	}
	// try from right to left
	n := len(a)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = a[n-1-i]
	}
	return canClearLeftRight(b)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r, k int
			fmt.Fscan(reader, &l, &r, &k)
			for i := l; i <= r; i++ {
				arr[i] += k
			}
		} else if t == 2 {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			seg := make([]int, r-l+1)
			copy(seg, arr[l:r+1])
			if canClear(seg) {
				fmt.Fprintln(writer, 1)
			} else {
				fmt.Fprintln(writer, 0)
			}
		}
	}
}
