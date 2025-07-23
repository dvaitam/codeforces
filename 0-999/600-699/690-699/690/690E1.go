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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	for ; q > 0; q-- {
		var h, w int
		if _, err := fmt.Fscan(reader, &h, &w); err != nil {
			return
		}
		img := make([][]int, h)
		for i := 0; i < h; i++ {
			row := make([]int, w)
			for j := 0; j < w; j++ {
				fmt.Fscan(reader, &row[j])
			}
			img[i] = row
		}
		diffs := make([]int64, h-1)
		for i := 0; i < h-1; i++ {
			var d int64
			r1 := img[i]
			r2 := img[i+1]
			for j := 0; j < w; j++ {
				diff := r1[j] - r2[j]
				if diff < 0 {
					diff = -diff
				}
				d += int64(diff)
			}
			diffs[i] = d
		}
		midIdx := h/2 - 1
		midVal := float64(diffs[midIdx])
		var sum float64
		for i := 0; i < len(diffs); i++ {
			if i == midIdx {
				continue
			}
			sum += float64(diffs[i])
		}
		avg := sum / float64(len(diffs)-1)
		if avg < 1 {
			avg = 1
		}
		if midVal > 1.3*avg {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
