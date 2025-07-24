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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	minR1 := int64(^uint64(0) >> 1) // max int64
	maxL1 := int64(0)
	for i := 0; i < n; i++ {
		var l, r int64
		fmt.Fscan(reader, &l, &r)
		if l > maxL1 {
			maxL1 = l
		}
		if r < minR1 {
			minR1 = r
		}
	}

	var m int
	fmt.Fscan(reader, &m)
	minR2 := int64(^uint64(0) >> 1)
	maxL2 := int64(0)
	for i := 0; i < m; i++ {
		var l, r int64
		fmt.Fscan(reader, &l, &r)
		if l > maxL2 {
			maxL2 = l
		}
		if r < minR2 {
			minR2 = r
		}
	}

	diff1 := maxL2 - minR1
	diff2 := maxL1 - minR2
	if diff1 < 0 {
		diff1 = 0
	}
	if diff2 < 0 {
		diff2 = 0
	}
	if diff1 > diff2 {
		fmt.Fprintln(writer, diff1)
	} else {
		fmt.Fprintln(writer, diff2)
	}
}
