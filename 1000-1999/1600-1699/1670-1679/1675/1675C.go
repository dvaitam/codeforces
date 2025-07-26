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
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	firstZero := n + 1
	lastOne := 0
	for i := 0; i < n; i++ {
		if s[i] == '0' && firstZero == n+1 {
			firstZero = i + 1
		}
		if s[i] == '1' {
			lastOne = i + 1
		}
	}
	start := lastOne
	if start == 0 {
		start = 1
	}
	end := firstZero
	if end == n+1 {
		end = n
	}
	if end >= start {
		fmt.Fprintln(writer, end-start+1)
	} else {
		fmt.Fprintln(writer, 0)
	}
}
