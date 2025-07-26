package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func main() {
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
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
	ops := 0
	for i := n - 2; i >= 0; i-- {
		for a[i] >= a[i+1] && a[i] > 0 {
			a[i] /= 2
			ops++
		}
		if a[i] >= a[i+1] {
			fmt.Fprintln(writer, -1)
			return
		}
	}
	fmt.Fprintln(writer, ops)
}
