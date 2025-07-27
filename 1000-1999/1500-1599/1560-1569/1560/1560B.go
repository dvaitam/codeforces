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
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		diff := abs(a - b)
		n := diff * 2
		if diff == 0 || max(a, b, c) > n {
			fmt.Fprintln(writer, -1)
			continue
		}
		half := diff
		var d int
		if c > half {
			d = c - half
		} else {
			d = c + half
		}
		fmt.Fprintln(writer, d)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(nums ...int) int {
	m := nums[0]
	for _, v := range nums[1:] {
		if v > m {
			m = v
		}
	}
	return m
}
