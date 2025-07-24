package main

import (
	"bufio"
	"fmt"
	"os"
)

func can(a []int, k int, D int) bool {
	mNeeded := (D + 1) / 2
	cnt := 0
	for _, v := range a {
		if v < mNeeded {
			cnt++
		}
	}
	if cnt > k {
		return false
	}
	rem := k - cnt
	l := 0
	bad := 0
	for r, v := range a {
		if v >= mNeeded && v < D {
			bad++
		}
		for bad > rem {
			y := a[l]
			if y >= mNeeded && y < D {
				bad--
			}
			l++
		}
		if r-l+1 >= 2 {
			return true
		}
	}
	return false
}

func solve(a []int, k int) int {
	lo, hi := 1, 1000000000
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if can(a, k, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		ans := solve(arr, k)
		fmt.Fprintln(writer, ans)
	}
}
