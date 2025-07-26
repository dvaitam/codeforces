package main

import (
	"bufio"
	"fmt"
	"os"
)

var fib [45]int64

func possible(n int, x, y int64) bool {
	for n > 1 {
		if fib[n-1] < y && y <= fib[n] {
			return false
		}
		if y <= fib[n-1] {
			x, y = y, fib[n]-x+1
		} else {
			y -= fib[n]
			x, y = y, fib[n]-x+1
		}
		n--
	}
	return true
}

func main() {
	fib[0], fib[1] = 1, 1
	for i := 2; i < len(fib); i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var x, y int64
		fmt.Fscan(reader, &n, &x, &y)
		if possible(n, x, y) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
