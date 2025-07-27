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
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	a := make([]int64, n)
	var prefixMax int64
	for i := 0; i < n; i++ {
		a[i] = b[i] + prefixMax
		if a[i] > prefixMax {
			prefixMax = a[i]
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, a[i])
	}
	writer.WriteByte('\n')
}
