package main

import (
	"bufio"
	"fmt"
	"os"
)

func isHillOrValley(a []int, i int) bool {
	if i <= 0 || i >= len(a)-1 {
		return false
	}
	return (a[i] > a[i-1] && a[i] > a[i+1]) || (a[i] < a[i-1] && a[i] < a[i+1])
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if n <= 2 {
			fmt.Fprintln(writer, 0)
			continue
		}
		hv := make([]bool, n)
		total := 0
		for i := 1; i < n-1; i++ {
			hv[i] = isHillOrValley(a, i)
			if hv[i] {
				total++
			}
		}
		ans := total
		for i := 0; i < n; i++ {
			orig := a[i]
			before := 0
			for j := i - 1; j <= i+1; j++ {
				if j > 0 && j < n-1 && hv[j] {
					before++
				}
			}
			if i > 0 {
				a[i] = a[i-1]
				after := 0
				for j := i - 1; j <= i+1; j++ {
					if j > 0 && j < n-1 && isHillOrValley(a, j) {
						after++
					}
				}
				if total-before+after < ans {
					ans = total - before + after
				}
			}
			if i < n-1 {
				a[i] = a[i+1]
				after := 0
				for j := i - 1; j <= i+1; j++ {
					if j > 0 && j < n-1 && isHillOrValley(a, j) {
						after++
					}
				}
				if total-before+after < ans {
					ans = total - before + after
				}
			}
			a[i] = orig
		}
		fmt.Fprintln(writer, ans)
	}
}
