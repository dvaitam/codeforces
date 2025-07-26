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

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	arr := make([]int64, n)
	times := make([]int, n)
	var sum int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		sum += arr[i]
		times[i] = 0
	}

	curTime := 0
	globalTime := -1
	var globalVal int64

	for ; q > 0; q-- {
		curTime++
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var idx int
			var x int64
			fmt.Fscan(reader, &idx, &x)
			idx--
			var old int64
			if times[idx] > globalTime {
				old = arr[idx]
			} else {
				old = globalVal
			}
			sum += x - old
			arr[idx] = x
			times[idx] = curTime
			fmt.Fprintln(writer, sum)
		} else {
			var x int64
			fmt.Fscan(reader, &x)
			globalVal = x
			globalTime = curTime
			sum = x * int64(n)
			fmt.Fprintln(writer, sum)
		}
	}
}
