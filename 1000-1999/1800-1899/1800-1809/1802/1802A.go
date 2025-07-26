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
		positives := 0
		negatives := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x > 0 {
				positives++
			} else {
				negatives++
			}
		}
		maxVals := make([]int, n)
		for i := 1; i <= n; i++ {
			if i <= positives {
				maxVals[i-1] = i
			} else {
				maxVals[i-1] = 2*positives - i
			}
		}
		minVals := make([]int, n)
		paired := negatives * 2
		for i := 1; i <= n; i++ {
			if i <= paired {
				if i%2 == 1 {
					minVals[i-1] = 1
				} else {
					minVals[i-1] = 0
				}
			} else {
				minVals[i-1] = i - paired
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, maxVals[i])
		}
		fmt.Fprintln(out)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, minVals[i])
		}
		fmt.Fprintln(out)
	}
}
