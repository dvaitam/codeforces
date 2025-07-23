package main

import (
	"bufio"
	"fmt"
	"os"
)

func sign(v int) int {
	if v > 0 {
		return 1
	}
	if v < 0 {
		return -1
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	fmt.Fscan(reader, &s)

	x, y := 0, 0
	last := 0
	coins := 0

	for _, c := range s {
		prev := sign(x - y)
		if c == 'U' {
			y++
		} else {
			x++
		}
		cur := sign(x - y)
		if cur == 0 {
			last = prev
		} else if prev == 0 {
			if last != 0 && cur != last {
				coins++
			}
		}
	}
	fmt.Fprintln(writer, coins)
}
