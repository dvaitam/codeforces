package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var x, y, l, r int64
	if _, err := fmt.Fscan(reader, &x, &y, &l, &r); err != nil {
		return
	}

	xs := make([]int64, 0, 64)
	for v := int64(1); ; {
		xs = append(xs, v)
		if v > r/x {
			break
		}
		v *= x
	}

	ys := make([]int64, 0, 64)
	for v := int64(1); ; {
		ys = append(ys, v)
		if v > r/y {
			break
		}
		v *= y
	}

	valsMap := make(map[int64]struct{})
	for _, a := range xs {
		for _, b := range ys {
			s := a + b
			if s >= l && s <= r {
				valsMap[s] = struct{}{}
			}
		}
	}

	vals := make([]int64, 0, len(valsMap)+2)
	for v := range valsMap {
		vals = append(vals, v)
	}
	vals = append(vals, l-1, r+1)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

	maxLen := int64(0)
	for i := 0; i < len(vals)-1; i++ {
		length := vals[i+1] - vals[i] - 1
		if length > maxLen {
			maxLen = length
		}
	}

	fmt.Fprintln(writer, maxLen)
}
