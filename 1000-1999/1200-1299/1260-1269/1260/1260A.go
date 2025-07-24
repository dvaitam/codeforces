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
	for i := 0; i < n; i++ {
		var c, sum int64
		fmt.Fscan(in, &c, &sum)
		k := c
		if k > sum {
			k = sum
		}
		q := sum / k
		r := sum % k
		cost := r*(q+1)*(q+1) + (k-r)*q*q
		fmt.Fprintln(out, cost)
	}
}
