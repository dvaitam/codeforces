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
	neg := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] < 0 {
			a[i] = -a[i]
			neg++
		}
	}
	for i := 0; i < neg; i++ {
		a[i] = -a[i]
	}
	ok := true
	for i := 1; i < n; i++ {
		if a[i-1] > a[i] {
			ok = false
			break
		}
	}
	if ok {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
