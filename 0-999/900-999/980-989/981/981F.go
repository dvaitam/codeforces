package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func feasible(a, b []int64, L int64, r int64) bool {
	n := len(a)
	// extended b with -L and +L
	ext := make([]int64, 3*n)
	for i := 0; i < n; i++ {
		ext[i] = b[i] - L
		ext[i+n] = b[i]
		ext[i+2*n] = b[i] + L
	}
	j1, j2 := 0, 0
	left, right := int64(-1<<60), int64(1<<60)
	for i := 0; i < n; i++ {
		ai := a[i]
		// move j1 to first >= ai-r
		for j1 < len(ext) && ext[j1] < ai-r {
			j1++
		}
		// move j2 to first > ai+r
		for j2 < len(ext) && ext[j2] <= ai+r {
			j2++
		}
		if j1 >= j2 {
			return false
		}
		low := int64(j1 - i)
		high := int64(j2 - 1 - i)
		if low > left {
			left = low
		}
		if high < right {
			right = high
		}
		if left > right {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var L int64
	if _, err := fmt.Fscan(reader, &n, &L); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	low, high := int64(0), L
	for low < high {
		mid := (low + high) / 2
		if feasible(a, b, L, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	fmt.Fprintln(writer, low)
}
