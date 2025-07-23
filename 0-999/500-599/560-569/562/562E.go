package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var x, y int64
	fmt.Fscan(in, &x, &y)
	var maxA, maxB int64
	for i := 0; i < n; i++ {
		var a, b int64
		fmt.Fscan(in, &a, &b)
		if a > maxA {
			maxA = a
		}
		if b > maxB {
			maxB = b
		}
	}
	var maxC, maxD int64
	for j := 0; j < m; j++ {
		var c, d int64
		fmt.Fscan(in, &c, &d)
		if c > maxC {
			maxC = c
		}
		if d > maxD {
			maxD = d
		}
	}
	if maxC > maxA && maxD > maxB {
		fmt.Println("Min")
	} else {
		fmt.Println("Max")
	}
}
