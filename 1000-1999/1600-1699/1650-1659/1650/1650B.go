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
		var l, r, a int
		fmt.Fscan(reader, &l, &r, &a)

		// Option 1: take x = r
		ans := r/a + r%a

		// Option 2: the largest number < r with remainder a-1, if within [l,r]
		candidateX := r - r%a - 1
		if candidateX >= l {
			candidate := candidateX/a + candidateX%a
			if candidate > ans {
				ans = candidate
			}
		}

		fmt.Fprintln(writer, ans)
	}
}
