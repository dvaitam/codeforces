package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit = 100000000

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func query(x, y int) int {
	fmt.Fprintf(writer, "? %d %d\n", x, y)
	writer.Flush()
	var d int
	fmt.Fscan(reader, &d)
	return d
}

func main() {
	// find y not on any horizontal line
	ySafe := limit
	for i := 0; i <= 10000; i++ {
		y := limit - i
		if query(0, y) != 0 {
			ySafe = y
			break
		}
	}

	vertical := make([]int, 0, 10000)
	x := -limit
	for x <= limit {
		d := query(x, ySafe)
		if d == 0 {
			vertical = append(vertical, x)
			x++
			continue
		}
		cand := x + d
		if cand <= limit && query(cand, ySafe) == 0 {
			vertical = append(vertical, cand)
			x = cand + 1
		} else {
			x = cand + 1
		}
	}

	// find x not on any vertical line
	xSafe := limit
	for i := 0; i <= 10000; i++ {
		x0 := limit - i
		if query(x0, 0) != 0 {
			xSafe = x0
			break
		}
	}

	horizontal := make([]int, 0, 10000)
	y := -limit
	for y <= limit {
		d := query(xSafe, y)
		if d == 0 {
			horizontal = append(horizontal, y)
			y++
			continue
		}
		cand := y + d
		if cand <= limit && query(xSafe, cand) == 0 {
			horizontal = append(horizontal, cand)
			y = cand + 1
		} else {
			y = cand + 1
		}
	}

	fmt.Fprintln(writer, "!")
	fmt.Fprintln(writer, len(vertical))
	for _, v := range vertical {
		fmt.Fprintln(writer, v)
	}
	fmt.Fprintln(writer, len(horizontal))
	for _, h := range horizontal {
		fmt.Fprintln(writer, h)
	}
	writer.Flush()
}
