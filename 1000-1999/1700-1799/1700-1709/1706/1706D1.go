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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		minA := 1<<31 - 1
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] < minA {
				minA = a[i]
			}
		}
		ans := 1<<31 - 1
		for m := 0; m <= minA; m++ {
			maxq := 0
			minq := 1<<31 - 1
			for _, ai := range a {
				var p int
				if m == 0 {
					p = k
				} else {
					p = ai / m
					if p > k {
						p = k
					}
					if p < 1 {
						p = 1
					}
				}
				q := ai / p
				if q > maxq {
					maxq = q
				}
				if q < minq {
					minq = q
				}
			}
			diff := maxq - minq
			if diff < ans {
				ans = diff
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
