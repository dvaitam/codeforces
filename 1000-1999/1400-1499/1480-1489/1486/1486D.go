package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(a []int, k int, mid int) bool {
	n := len(a)
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if a[i-1] >= mid {
			prefix[i] = prefix[i-1] + 1
		} else {
			prefix[i] = prefix[i-1] - 1
		}
	}
	minPrefix := 0
	for i := k; i <= n; i++ {
		if prefix[i]-minPrefix > 0 {
			return true
		}
		if prefix[i-k+1] < minPrefix {
			minPrefix = prefix[i-k+1]
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	fmt.Fscan(reader, &n, &k)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	low, high := 1, n
	res := 1
	for low <= high {
		mid := (low + high) / 2
		if check(a, k, mid) {
			res = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	fmt.Fprintln(writer, res)
}
