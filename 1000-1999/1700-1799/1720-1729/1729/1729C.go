package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	idx int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		startVal := int(s[0] - 'a')
		endVal := int(s[n-1] - 'a')
		minv, maxv := startVal, endVal
		if minv > maxv {
			minv, maxv = maxv, minv
		}
		var arr []pair
		for i := 1; i < n-1; i++ {
			v := int(s[i] - 'a')
			if v >= minv && v <= maxv {
				arr = append(arr, pair{v, i + 1})
			}
		}
		if startVal <= endVal {
			sort.Slice(arr, func(i, j int) bool {
				if arr[i].val == arr[j].val {
					return arr[i].idx < arr[j].idx
				}
				return arr[i].val < arr[j].val
			})
		} else {
			sort.Slice(arr, func(i, j int) bool {
				if arr[i].val == arr[j].val {
					return arr[i].idx < arr[j].idx
				}
				return arr[i].val > arr[j].val
			})
		}
		path := make([]int, 0, len(arr)+2)
		path = append(path, 1)
		for _, p := range arr {
			path = append(path, p.idx)
		}
		path = append(path, n)
		cost := abs(startVal - endVal)
		fmt.Fprintln(writer, cost, len(path))
		for i, pos := range path {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprintf(writer, "%d", pos)
		}
		writer.WriteByte('\n')
	}
}
