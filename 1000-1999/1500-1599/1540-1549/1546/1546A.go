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
		a := make([]int, n)
		b := make([]int, n)
		sumA, sumB := 0, 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sumA += a[i]
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
			sumB += b[i]
		}
		if sumA != sumB {
			fmt.Fprintln(out, -1)
			continue
		}
		from := []int{}
		to := []int{}
		for i := 0; i < n; i++ {
			diff := a[i] - b[i]
			for diff > 0 {
				from = append(from, i+1)
				diff--
			}
			for diff < 0 {
				to = append(to, i+1)
				diff++
			}
		}
		fmt.Fprintln(out, len(from))
		for i := 0; i < len(from); i++ {
			fmt.Fprintln(out, from[i], to[i])
		}
	}
}
