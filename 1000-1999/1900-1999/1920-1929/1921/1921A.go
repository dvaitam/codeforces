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
		minX, maxX := 1<<31-1, -1<<31
		minY, maxY := 1<<31-1, -1<<31
		for i := 0; i < 4; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
		side := maxX - minX // or maxY - minY, they should be equal
		area := side * side
		fmt.Fprintln(writer, area)
	}
}
