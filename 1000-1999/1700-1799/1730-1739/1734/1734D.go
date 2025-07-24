package main

import (
	"bufio"
	"fmt"
	"os"
)

func canEscape(a []int64, k int) bool {
	n := len(a) - 1 // 1-indexed
	l, r := k, k
	cur := a[k]
	for {
		if l == 1 || r == n {
			return true
		}
		moved := false
		// try expand left
		acc := int64(0)
		for l > 1 && cur+acc+a[l-1] >= 0 {
			l--
			acc += a[l]
			if acc > 0 {
				cur += acc
				acc = 0
				moved = true
			}
		}
		if l == 1 {
			return true
		}
		// try expand right
		acc = 0
		for r < n && cur+acc+a[r+1] >= 0 {
			r++
			acc += a[r]
			if acc > 0 {
				cur += acc
				acc = 0
				moved = true
			}
		}
		if r == n {
			return true
		}
		if !moved {
			return false
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		if canEscape(arr, k) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
