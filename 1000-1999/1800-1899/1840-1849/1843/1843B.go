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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		var sum int64
		for i := 0; i < n; i++ {
			if arr[i] < 0 {
				sum -= arr[i]
			} else {
				sum += arr[i]
			}
		}

		ops := 0
		i := 0
		for i < n {
			for i < n && arr[i] > 0 {
				i++
			}
			if i == n {
				break
			}
			hasNeg := false
			for i < n && arr[i] <= 0 {
				if arr[i] < 0 {
					hasNeg = true
				}
				i++
			}
			if hasNeg {
				ops++
			}
		}

		fmt.Fprintf(out, "%d %d\n", sum, ops)
	}
}
