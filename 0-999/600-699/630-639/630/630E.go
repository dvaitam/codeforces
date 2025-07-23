package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var x1, y1, x2, y2 int64
	if _, err := fmt.Fscan(reader, &x1, &y1, &x2, &y2); err != nil {
		return
	}
	width := x2 - x1 + 1
	height := y2 - y1 + 1
	total := width * height
	ans := total / 2
	if total%2 == 1 {
		// parity of x1 + y1
		if ((x1 ^ y1) & 1) == 0 {
			ans++
		}
	}
	fmt.Println(ans)
}
