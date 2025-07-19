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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	fmt.Fprintln(writer, n)
	for i := 0; i < n; i++ {
		minv := int(2147483647)
		minj := i
		for j := i; j < n; j++ {
			if a[j] < minv {
				minv = a[j]
				minj = j
			}
		}
		a[i], a[minj] = a[minj], a[i]
		fmt.Fprintln(writer, i, minj)
	}
}
