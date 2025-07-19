package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		total := n + 2
		arr := make([]int64, total)
		var sum int64
		for i := 0; i < total; i++ {
			fmt.Fscan(reader, &arr[i])
			sum += arr[i]
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		found := false
		var res []int64
		// Try largest as y
		y := arr[total-1]
		x := sum - 2*y
		if idx := indexOf(arr, 0, total-1, x); idx >= 0 {
			found = true
			for i, v := range arr {
				if i == idx || i == total-1 {
					continue
				}
				res = append(res, v)
			}
		} else {
			// Try second largest as y and largest as x
			y = arr[total-2]
			x = sum - 2*y
			if x == arr[total-1] {
				found = true
				for i, v := range arr {
					if i == total-2 || i == total-1 {
						continue
					}
					res = append(res, v)
				}
			}
		}
		if !found {
			fmt.Fprintln(writer, -1)
		} else {
			for i, v := range res {
				if i > 0 {
					writer.WriteByte(' ')
				}
				writer.WriteString(strconv.FormatInt(v, 10))
			}
			writer.WriteByte('\n')
		}
	}
}

// indexOf finds x in arr[lo:hi] (hi exclusive), returns index or -1
func indexOf(a []int64, lo, hi int, x int64) int {
	idx := sort.Search(hi-lo, func(i int) bool { return a[lo+i] >= x })
	if idx < hi-lo && a[lo+idx] == x {
		return lo + idx
	}
	return -1
}
