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
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		var g int64
		for i := 0; i < n/2; i++ {
			diff := arr[i] - arr[n-1-i]
			if diff < 0 {
				diff = -diff
			}
			g = gcd(g, diff)
		}
		if g == 0 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, g)
		}
	}
}
