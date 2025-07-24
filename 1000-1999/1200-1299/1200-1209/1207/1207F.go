package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxN = 500000
	B    = 710 // ~sqrt(maxN)
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	arr := make([]int, maxN+1)
	small := make([][]int, B)
	for i := 0; i < B; i++ {
		small[i] = make([]int, B)
	}

	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(reader, &t, &x, &y)
		if t == 1 {
			arr[x] += y
			for m := 1; m < B; m++ {
				small[m][x%m] += y
			}
		} else {
			if x < B {
				fmt.Fprintln(writer, small[x][y])
			} else {
				sum := 0
				start := y
				if start == 0 {
					start = x
				}
				for j := start; j <= maxN; j += x {
					sum += arr[j]
				}
				fmt.Fprintln(writer, sum)
			}
		}
	}
}
