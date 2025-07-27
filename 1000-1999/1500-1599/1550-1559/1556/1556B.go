package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func cost(indices []int, start int) int64 {
	var sum int64
	for i, pos := range indices {
		target := start + i*2
		sum += absInt64(int64(pos - target))
	}
	return sum
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
		var n int
		fmt.Fscan(reader, &n)
		evenIdx := make([]int, 0, n)
		oddIdx := make([]int, 0, n)
		for i := 1; i <= n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x%2 == 0 {
				evenIdx = append(evenIdx, i)
			} else {
				oddIdx = append(oddIdx, i)
			}
		}
		ce := len(evenIdx)
		co := len(oddIdx)
		if n%2 == 0 {
			if ce != co {
				fmt.Fprintln(writer, -1)
				continue
			}
			costEvenStart := cost(evenIdx, 1) // even at positions 1,3,...
			costOddStart := cost(oddIdx, 1)   // odd at positions 1,3,...
			if costEvenStart < costOddStart {
				fmt.Fprintln(writer, costEvenStart)
			} else {
				fmt.Fprintln(writer, costOddStart)
			}
		} else {
			if absInt64(int64(ce-co)) != 1 {
				fmt.Fprintln(writer, -1)
				continue
			}
			if ce > co {
				fmt.Fprintln(writer, cost(evenIdx, 1))
			} else {
				fmt.Fprintln(writer, cost(oddIdx, 1))
			}
		}
	}
}
