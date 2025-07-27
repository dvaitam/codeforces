package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var x1, y1, x2, y2 int64
	fmt.Fscan(reader, &x1, &y1)
	fmt.Fscan(reader, &x2, &y2)

	cross := x1*y2 - x2*y1
	if cross < 0 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		cross = -cross
	}

	if n != cross {
		fmt.Fprintln(writer, "NO")
		return
	}

	dx := cross / gcd(y1, y2)
	dy := cross / gcd(x1, x2)

	fmt.Fprintln(writer, "YES")
	for i := int64(0); i < dx && n > 0; i++ {
		for j := int64(0); j < dy && n > 0; j++ {
			fmt.Fprintln(writer, i, j)
			n--
		}
	}
}
