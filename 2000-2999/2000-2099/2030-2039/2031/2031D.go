package main

import (
	"bufio"
	"fmt"
	"os"
)

const infValue = int(1e9 + 7)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		suffixMin := make([]int, n+2)
		suffixMin[n+1] = infValue
		suffixMin[n] = a[n]
		for i := n - 1; i >= 1; i-- {
			if a[i] < suffixMin[i+1] {
				suffixMin[i] = a[i]
			} else {
				suffixMin[i] = suffixMin[i+1]
			}
		}

		ans := make([]int, n+1)
		start := 1
		for start <= n {
			maxVal := a[start]
			j := start
			for {
				if j == n || maxVal <= suffixMin[j+1] {
					for k := start; k <= j; k++ {
						ans[k] = maxVal
					}
					start = j + 1
					break
				}
				j++
				if a[j] > maxVal {
					maxVal = a[j]
				}
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
