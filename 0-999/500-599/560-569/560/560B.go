package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func canFit(A, B, a2, b2, a3, b3 int) bool {
	for r1 := 0; r1 < 2; r1++ {
		x1, y1 := a2, b2
		if r1 == 1 {
			x1, y1 = b2, a2
		}
		for r2 := 0; r2 < 2; r2++ {
			x2, y2 := a3, b3
			if r2 == 1 {
				x2, y2 = b3, a3
			}
			if x1+x2 <= A && max(y1, y2) <= B {
				return true
			}
			if y1+y2 <= B && max(x1, x2) <= A {
				return true
			}
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var a1, b1, a2, b2, a3, b3 int
	fmt.Fscan(reader, &a1, &b1)
	fmt.Fscan(reader, &a2, &b2)
	fmt.Fscan(reader, &a3, &b3)

	if canFit(a1, b1, a2, b2, a3, b3) || canFit(b1, a1, a2, b2, a3, b3) {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
