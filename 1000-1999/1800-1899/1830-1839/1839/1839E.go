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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	sum := 0
	maxv := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		sum += a[i]
		if a[i] > maxv {
			maxv = a[i]
		}
	}

	if sum%2 == 1 || maxv > sum-maxv {
		fmt.Fprintln(out, "First")
	} else {
		fmt.Fprintln(out, "Second")
	}
}
