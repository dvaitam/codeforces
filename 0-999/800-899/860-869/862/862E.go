package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func minAbs(arr []int64, x int64) int64 {
	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= x })
	best := int64(1<<63 - 1)
	if i < len(arr) {
		diff := arr[i] - x
		if diff < 0 {
			diff = -diff
		}
		if diff < best {
			best = diff
		}
	}
	if i > 0 {
		diff := x - arr[i-1]
		if diff < 0 {
			diff = -diff
		}
		if diff < best {
			best = diff
		}
	}
	return best
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	A := int64(0)
	for i := 1; i <= n; i++ {
		var v int64
		fmt.Fscan(in, &v)
		if i%2 == 1 {
			A += v
		} else {
			A -= v
		}
	}
	bPrefix := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		var v int64
		fmt.Fscan(in, &v)
		if i%2 == 1 {
			bPrefix[i] = bPrefix[i-1] + v
		} else {
			bPrefix[i] = bPrefix[i-1] - v
		}
	}
	arr := make([]int64, m-n+1)
	for j := 0; j <= m-n; j++ {
		val := bPrefix[j+n] - bPrefix[j]
		if j%2 == 1 {
			val = -val
		}
		arr[j] = val
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	fmt.Fprintln(out, minAbs(arr, A))
	for ; q > 0; q-- {
		var l, r int
		var x int64
		fmt.Fscan(in, &l, &r, &x)
		if (r-l+1)%2 == 1 {
			if l%2 == 1 {
				A += x
			} else {
				A -= x
			}
		}
		fmt.Fprintln(out, minAbs(arr, A))
	}
}
