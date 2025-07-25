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
		negCount := 0
		zero := false
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if arr[i] == 0 {
				zero = true
			} else if arr[i] < 0 {
				negCount++
			}
		}
		if zero || negCount%2 == 1 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, 1, 0)
		}
	}
}
