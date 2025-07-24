package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, l, r int
	if _, err := fmt.Fscan(in, &n, &l, &r); err != nil {
		return
	}

	minSum := int64(n - l)
	minSum += int64(1<<l) - 1

	maxSum := int64(1<<r) - 1
	maxSum += int64(n-r) * int64(1<<(r-1))

	fmt.Fprintln(out, minSum, maxSum)
}
