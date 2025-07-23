package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a, b, c, d int
	if _, err := fmt.Fscan(in, &n, &a, &b, &c, &d); err != nil {
		return
	}

	low := max4(1, 1-(b-c), 1-(a-d), 1-(a+b-c-d))
	high := min4(n, n-(b-c), n-(a-d), n-(a+b-c-d))
	if low > high {
		fmt.Println(0)
		return
	}
	ans := (high - low + 1) * n
	fmt.Println(ans)
}

func max4(a, b, c, d int) int {
	if a < b {
		a = b
	}
	if a < c {
		a = c
	}
	if a < d {
		a = d
	}
	return a
}

func min4(a, b, c, d int) int {
	if a > b {
		a = b
	}
	if a > c {
		a = c
	}
	if a > d {
		a = d
	}
	return a
}
