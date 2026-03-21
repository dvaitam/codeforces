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

	for i := 0; i < t; i++ {
		var x, y int64
		fmt.Fscan(reader, &x, &y)

		A := x + y + 1
		B := y - x

		var digits []int
		for A != 1 || B != 0 {
			var d int
			var remA, remB int64
			for d = 0; d < 4; d++ {
				if d == 0 {
					remA, remB = A-1, B
				} else if d == 1 {
					remA, remB = A, B-1
				} else if d == 2 {
					remA, remB = A+1, B
				} else {
					remA, remB = A, B+1
				}
				if remA%2 == 0 && (remA+remB)%4 == 0 {
					break
				}
			}
			digits = append(digits, d)
			A, B = remB/2, -remA/2
		}

		var n int64 = 0
		for j := len(digits) - 1; j >= 0; j-- {
			d := digits[j]
			shift := (n/4) % 2
			c := (d + int(shift)) % 4
			n = n*4 + int64(c)
		}

		curve := (n % 4) + 1
		pos := (n / 4) + 1
		fmt.Fprintf(writer, "%d %d\n", curve, pos)
	}
}
