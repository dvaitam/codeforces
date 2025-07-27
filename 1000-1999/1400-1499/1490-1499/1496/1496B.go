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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		set := make(map[int64]struct{}, n)
		maxVal := int64(0)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			set[x] = struct{}{}
			if i == 0 || x > maxVal {
				maxVal = x
			}
		}
		if k == 0 {
			fmt.Fprintln(writer, len(set))
			continue
		}
		// compute mex
		mex := int64(0)
		for {
			if _, ok := set[mex]; !ok {
				break
			}
			mex++
		}
		if mex > maxVal {
			// we can add k new unique elements starting from maxVal+1
			fmt.Fprintln(writer, int64(len(set))+k)
			continue
		}
		// candidate value to insert
		y := (mex + maxVal + 1) / 2
		if _, ok := set[y]; !ok {
			fmt.Fprintln(writer, len(set)+1)
		} else {
			fmt.Fprintln(writer, len(set))
		}
	}
}
