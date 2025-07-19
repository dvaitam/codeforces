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

	var tt, n, x, y int
	if _, err := fmt.Fscan(reader, &tt); err != nil {
		return
	}
	for tt > 0 {
		tt--
		fmt.Fscan(reader, &n, &x, &y)
		// invalid if both x and y positive
		if (x > 0 && y > 0) || (x == 0 && y == 0) {
			fmt.Fprintln(writer, -1)
			continue
		}
		if x == 0 {
			x, y = y, x
		}
		if x == 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		if (n-1)%x != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		a, b, cnt := 1, 2, 0
		for i := 0; i < n-1; i++ {
			if cnt < x {
				fmt.Fprint(writer, a, " ")
				b++
				cnt++
			} else {
				fmt.Fprint(writer, b, " ")
				a = b + 1
				// swap a, b
				a, b = b, a
				cnt = 1
			}
		}
		fmt.Fprintln(writer)
	}
}
