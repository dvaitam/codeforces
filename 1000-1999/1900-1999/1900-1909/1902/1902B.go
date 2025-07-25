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

	var tc int
	fmt.Fscan(reader, &tc)
	for ; tc > 0; tc-- {
		var n, P, l, t int64
		fmt.Fscan(reader, &n, &P, &l, &t)

		tasks := (n + 6) / 7

		low, high := int64(0), n
		for low < high {
			mid := (low + high) / 2
			taskDone := 2 * mid
			if taskDone > tasks {
				taskDone = tasks
			}
			points := mid*l + taskDone*t
			if points >= P {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(writer, n-low)
	}
}
