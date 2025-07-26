package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isMagic(a, b, x int64) bool {
	g := gcd(a, b)
	if x%g != 0 {
		return false
	}
	a /= g
	b /= g
	x /= g
	for a != 0 && b != 0 && max64(a, b) >= x {
		if a == x || b == x {
			return true
		}
		if a < b {
			a, b = b, a
		}
		if (a-x)%b == 0 {
			return true
		}
		a %= b
	}
	return a == x || b == x
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b, x int64
		fmt.Fscan(reader, &a, &b, &x)
		if isMagic(a, b, x) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
