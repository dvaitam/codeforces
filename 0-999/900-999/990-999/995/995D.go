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

	var n, q int
	fmt.Fscan(reader, &n, &q)
	size := 1 << uint(n)
	a := make([]int64, size)
	var sum int64
	for i := 0; i < size; i++ {
		fmt.Fscan(reader, &a[i])
		sum += a[i]
	}
	denom := float64(size)
	fmt.Fprintf(writer, "%.7f\n", float64(sum)/denom)
	for i := 0; i < q; i++ {
		var x int
		var y int64
		fmt.Fscan(reader, &x, &y)
		sum += y - a[x]
		a[x] = y
		fmt.Fprintf(writer, "%.7f\n", float64(sum)/denom)
	}
}
