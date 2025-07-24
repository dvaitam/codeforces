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
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if n == 0 {
		fmt.Fprintln(out, 0)
		return
	}
	minVal, maxVal := a[0], a[0]
	for _, v := range a {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	count := 0
	for _, v := range a {
		if v > minVal && v < maxVal {
			count++
		}
	}
	fmt.Fprintln(out, count)
}
