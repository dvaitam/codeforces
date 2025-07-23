package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var x1, y1, x2, y2 int64
	if _, err := fmt.Fscan(in, &x1, &y1, &x2, &y2); err != nil {
		// attempt scanning line by line
		if _, err2 := fmt.Fscan(in, &x2, &y2); err2 != nil {
			return
		}
	}
	dx := abs64(x1 - x2)
	dy := abs64(y1 - y2)
	if dx > dy {
		fmt.Println(dx)
	} else {
		fmt.Println(dy)
	}
}
