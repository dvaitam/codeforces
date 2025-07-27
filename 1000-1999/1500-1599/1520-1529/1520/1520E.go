package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
		var s string
		fmt.Fscan(reader, &s)

		positions := make([]int, 0)
		for i, ch := range s {
			if ch == '*' {
				positions = append(positions, i)
			}
		}
		m := len(positions)
		if m == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		mid := m / 2
		base := positions[mid] - mid
		var res int64
		for i, p := range positions {
			res += int64(abs(p - (base + i)))
		}
		fmt.Fprintln(writer, res)
	}
}
