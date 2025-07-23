package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)

	x := make([]int, n+1)
	y := make([]int, n+1)
	for i := 0; i < n; i++ {
		x[i+1] = x[i]
		y[i+1] = y[i]
		switch s[i] {
		case 'U':
			y[i+1]++
		case 'D':
			y[i+1]--
		case 'L':
			x[i+1]--
		case 'R':
			x[i+1]++
		}
	}

	cnt := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			if x[j]-x[i] == 0 && y[j]-y[i] == 0 {
				cnt++
			}
		}
	}

	fmt.Fprintln(writer, cnt)
}
