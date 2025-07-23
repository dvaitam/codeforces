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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	var maxSum int64
	for l := 0; l < n; l++ {
		var orA, orB int64
		for r := l; r < n; r++ {
			orA |= a[r]
			orB |= b[r]
			sum := orA + orB
			if sum > maxSum {
				maxSum = sum
			}
		}
	}

	fmt.Fprintln(writer, maxSum)
}
