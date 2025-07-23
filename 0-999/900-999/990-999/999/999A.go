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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	left := 0
	for left < n && a[left] <= k {
		left++
	}

	right := n - 1
	for right >= left && a[right] <= k {
		right--
	}

	fmt.Fprintln(writer, left+(n-1-right))
}
