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

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve() {
	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)
	prefix0 := make([]int, n+1)
	prefix1 := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix0[i] = prefix0[i-1]
		prefix1[i] = prefix1[i-1]
		if s[i-1] == '0' {
			prefix0[i]++
		} else {
			prefix1[i]++
		}
	}
	total1 := prefix1[n]
	bestI := 0
	bestD := n*2 + 5
	for i := 0; i <= n; i++ {
		zerosLeft := prefix0[i]
		if zerosLeft < (i+1)/2 {
			continue
		}
		onesRight := total1 - prefix1[i]
		if onesRight < (n-i+1)/2 {
			continue
		}
		d := absInt(n - 2*i)
		if d < bestD {
			bestD = d
			bestI = i
		}
	}
	fmt.Fprintln(writer, bestI)
}

func main() {
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve()
	}
}
