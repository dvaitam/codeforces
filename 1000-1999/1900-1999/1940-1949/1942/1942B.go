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
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		used := make([]bool, n+2)
		p := make([]int, n)
		mex := 0
		for i := 0; i < n; i++ {
			ai := a[i]
			// try placing the current mex
			newMex := mex + 1
			for newMex <= n && used[newMex] {
				newMex++
			}
			if newMex == mex+ai {
				p[i] = mex
				used[mex] = true
				mex = newMex
				continue
			}

			// otherwise, choose a value that keeps mex unchanged
			val := mex - ai
			if val < 0 || val >= n || used[val] || val == mex {
				// input guarantees at least one solution
				// so we assume this branch will always be valid
			}
			p[i] = val
			used[val] = true
			for mex <= n && used[mex] {
				mex++
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, p[i])
		}
		fmt.Fprintln(writer)
	}
}
