package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt of contest 1690B (Array Decrements).
// We need to check if there exists a non-negative integer K such that for
// every i, b[i] = max(0, a[i]-K). For indices where b[i] > 0, K must equal
// a[i]-b[i]. For indices where b[i] == 0, a[i] must not exceed K.
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
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		diff := -1
		ok := true
		for i := 0; i < n; i++ {
			if b[i] > a[i] {
				ok = false
				break
			}
			if b[i] > 0 {
				d := a[i] - b[i]
				if diff == -1 {
					diff = d
				} else if diff != d {
					ok = false
					break
				}
			}
		}
		if ok && diff != -1 {
			for i := 0; i < n; i++ {
				if b[i] == 0 && a[i] > diff {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
