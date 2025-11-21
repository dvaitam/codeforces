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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		countZero := 0
		allZero := true
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] == 0 {
				countZero++
			} else {
				allZero = false
			}
		}
		if allZero {
			fmt.Fprintln(out, 0)
			continue
		}
		if countZero == 0 {
			fmt.Fprintln(out, 1)
			continue
		}
		left, right := 0, n-1
		for left < n && a[left] == 0 {
			left++
		}
		for right >= 0 && a[right] == 0 {
			right--
		}
		hasZeroInside := false
		if left <= right {
			for i := left; i <= right; i++ {
				if a[i] == 0 {
					hasZeroInside = true
					break
				}
			}
		}
		if hasZeroInside {
			fmt.Fprintln(out, 2)
		} else {
			fmt.Fprintln(out, 1)
		}
	}
}
