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

	var n, k, x int
	if _, err := fmt.Fscan(reader, &n, &k, &x); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	sum := 0
	// chores are sorted nondecreasing, replace last k with x
	for i := 0; i < n-k; i++ {
		sum += a[i]
	}
	sum += k * x

	fmt.Fprintln(writer, sum)
}
