package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var xp, yp, xv, yv int
	if _, err := fmt.Fscan(in, &xp, &yp, &xv, &yv); err != nil {
		return
	}

	if (xp <= xv && yp <= yv) || (xp+yp <= max(xv, yv)) {
		fmt.Println("Polycarp")
	} else {
		fmt.Println("Vasiliy")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
