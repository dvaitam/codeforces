package main

import (
	"bufio"
	"fmt"
	"os"
)

func canTransform(a, b []int) bool {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] > b[i] {
			return false
		}
	}
	for i := 0; i < n; i++ {
		x := b[i]
		l := i
		for l > 0 && a[l-1] <= x && b[l-1] >= x {
			l--
		}
		r := i
		for r < n-1 && a[r+1] <= x && b[r+1] >= x {
			r++
		}
		found := false
		for j := l; j <= r; j++ {
			if a[j] == x {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		if canTransform(a, b) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
