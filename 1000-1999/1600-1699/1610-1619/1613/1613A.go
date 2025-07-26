package main

import (
	"bufio"
	"fmt"
	"os"
)

func digits(x int64) int {
	d := 0
	for x > 0 {
		d++
		x /= 10
	}
	return d
}

func pow10(k int) int64 {
	res := int64(1)
	for i := 0; i < k; i++ {
		res *= 10
	}
	return res
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
		var x1, p1 int64
		var x2, p2 int64
		fmt.Fscan(reader, &x1, &p1)
		fmt.Fscan(reader, &x2, &p2)

		d1 := digits(x1)
		d2 := digits(x2)

		len1 := d1 + int(p1)
		len2 := d2 + int(p2)

		if len1 > len2 {
			fmt.Fprintln(writer, ">")
			continue
		}
		if len1 < len2 {
			fmt.Fprintln(writer, "<")
			continue
		}
		diff := d1 - d2
		if diff > 0 {
			x2 *= pow10(diff)
			p2 -= int64(diff)
		} else if diff < 0 {
			x1 *= pow10(-diff)
			p1 += int64(diff)
		}
		if x1 > x2 {
			fmt.Fprintln(writer, ">")
		} else if x1 < x2 {
			fmt.Fprintln(writer, "<")
		} else {
			fmt.Fprintln(writer, "=")
		}
	}
}
