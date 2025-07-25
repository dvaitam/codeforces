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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		var year int64 = 1
		var signYear int64
		for i := 0; i < n; i++ {
			ai := a[i]
			if year%ai == 0 {
				signYear = year
			} else {
				signYear = (year + ai - 1) / ai * ai
			}
			year = signYear + 1
		}
		fmt.Fprintln(out, signYear)
	}
}
