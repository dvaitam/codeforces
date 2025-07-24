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
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		var negCount int
		hasZero := false
		var sumAbs int64
		var minAbs int64 = 1<<63 - 1
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			if arr[i] < 0 {
				negCount++
			}
			if arr[i] == 0 {
				hasZero = true
			}
			v := arr[i]
			if v < 0 {
				v = -v
			}
			sumAbs += v
			if v < minAbs {
				minAbs = v
			}
		}
		if negCount%2 == 1 && !hasZero {
			sumAbs -= 2 * minAbs
		}
		fmt.Fprintln(writer, sumAbs)
	}
}
