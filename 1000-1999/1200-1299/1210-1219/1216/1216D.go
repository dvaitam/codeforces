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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int64, n)
	var maxA int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		if arr[i] > maxA {
			maxA = arr[i]
		}
	}

	var g int64
	for i := 0; i < n; i++ {
		diff := maxA - arr[i]
		g = gcd(g, diff)
	}

	var y int64
	for i := 0; i < n; i++ {
		y += (maxA - arr[i]) / g
	}

	fmt.Fprintln(writer, y, g)
}
