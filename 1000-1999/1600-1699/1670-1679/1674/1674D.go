package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func main() {
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve()
	}
}

func solve() {
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	start := n % 2
	for i := start; i+1 < n; i += 2 {
		if a[i] > a[i+1] {
			a[i], a[i+1] = a[i+1], a[i]
		}
	}
	for i := 1; i < n; i++ {
		if a[i-1] > a[i] {
			fmt.Fprintln(writer, "NO")
			return
		}
	}
	fmt.Fprintln(writer, "YES")
}
