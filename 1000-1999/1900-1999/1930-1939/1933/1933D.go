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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		// find minimum and its count
		minVal := arr[0]
		for _, v := range arr {
			if v < minVal {
				minVal = v
			}
		}
		cntMin := 0
		for _, v := range arr {
			if v == minVal {
				cntMin++
			}
		}
		possible := false
		if cntMin == 1 {
			possible = true
		} else {
			for _, v := range arr {
				if v%minVal != 0 {
					possible = true
					break
				}
			}
		}
		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
