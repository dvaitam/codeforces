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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		r := int(k % int64(n+1))
		freq := make([]int, n+1)
		for _, v := range a {
			freq[v] = 1
		}
		mex := 0
		for mex <= n && freq[mex] == 1 {
			mex++
		}
		x := make([]int, r+1)
		for step := 1; step <= r; step++ {
			x[step] = mex
			leave := a[n-step]
			freq[leave] = 0
			freq[mex] = 1
			if leave < mex {
				mex = leave
			} else {
				for mex <= n && freq[mex] == 1 {
					mex++
				}
			}
		}
		res := make([]int, n)
		for i := 0; i < r; i++ {
			res[i] = x[r-i]
		}
		for i := r; i < n; i++ {
			res[i] = a[i-r]
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, res[i])
		}
		writer.WriteByte('\n')
	}
}
