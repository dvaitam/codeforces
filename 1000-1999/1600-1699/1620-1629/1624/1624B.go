package main

import (
	"bufio"
	"fmt"
	"os"
)

func canMakeAP(a, b, c int64) bool {
	target := 2*b - c
	if target > 0 && target%a == 0 {
		return true
	}
	sum := a + c
	if sum%2 == 0 {
		mid := sum / 2
		if mid > 0 && mid%b == 0 {
			return true
		}
	}
	target = 2*b - a
	if target > 0 && target%c == 0 {
		return true
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b, c int64
		fmt.Fscan(in, &a, &b, &c)
		if canMakeAP(a, b, c) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
