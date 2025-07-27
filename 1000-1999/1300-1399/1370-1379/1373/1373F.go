package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		if possible(a, b) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

func simulate(a, b []int64, x int64) (int64, bool) {
	prev := x
	for i := 0; i < len(a); i++ {
		use := prev
		if use > a[i] {
			use = a[i]
		}
		need := a[i] - use
		if need > b[i] {
			return 0, false
		}
		prev = b[i] - need
	}
	return prev, true
}

func possible(a, b []int64) bool {
	low := int64(0)
	high := b[len(b)-1]
	for low <= high {
		mid := (low + high) / 2
		val, ok := simulate(a, b, mid)
		if !ok {
			low = mid + 1
			continue
		}
		if val >= mid {
			return true
		}
		high = mid - 1
	}
	return false
}
