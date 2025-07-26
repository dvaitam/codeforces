package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPowerOfTwo(x int) bool {
	return x > 0 && x&(x-1) == 0
}

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
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		ok := true
		for i := 0; i < n-1; i++ {
			if !isPowerOfTwo(i+1) && a[i] > a[i+1] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
