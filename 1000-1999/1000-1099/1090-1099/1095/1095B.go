package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)
	// Removing one element: either remove the largest or the smallest
	// Option1: remove largest -> a[n-2] - a[0]
	// Option2: remove smallest -> a[n-1] - a[1]
	op1 := a[n-2] - a[0]
	op2 := a[n-1] - a[1]
	if op1 < op2 {
		fmt.Fprintln(writer, op1)
	} else {
		fmt.Fprintln(writer, op2)
	}
}
