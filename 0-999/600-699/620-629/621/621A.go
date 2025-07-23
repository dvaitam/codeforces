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
	var sum int64
	const inf = int64(1<<63 - 1)
	minOdd := inf
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		sum += x
		if x%2 != 0 && x < minOdd {
			minOdd = x
		}
	}
	if sum%2 != 0 {
		if minOdd != inf {
			sum -= minOdd
		}
	}
	fmt.Fprintln(writer, sum)
}
