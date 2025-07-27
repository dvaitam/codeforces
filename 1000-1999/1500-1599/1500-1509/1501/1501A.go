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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i], &b[i])
		}
		tm := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &tm[i])
		}
		prevDepart := 0
		prevB := 0
		for i := 0; i < n; i++ {
			travel := a[i] - prevB + tm[i]
			arrival := prevDepart + travel
			if i == n-1 {
				fmt.Fprintln(out, arrival)
				break
			}
			stay := (b[i] - a[i] + 1) / 2
			depart := arrival + stay
			if depart < b[i] {
				depart = b[i]
			}
			prevDepart = depart
			prevB = b[i]
		}
	}
}
