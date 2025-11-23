package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n int
		var x, y int64
		fmt.Fscan(in, &n, &x, &y)

		var s string
		fmt.Fscan(in, &s)

		var k4, k8 int64
		for _, c := range s {
			if c == '4' {
				k4++
			} else {
				k8++
			}
		}

		ax := abs(x)
		ay := abs(y)

		mdist := ax + ay
		cdist := ax
		if ay > cdist {
			cdist = ay
		}

		if mdist <= 2*k8+k4 && cdist <= k4+k8 {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}
}

