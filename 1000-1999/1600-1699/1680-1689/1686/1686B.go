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

	var t, n int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		res := 0
		for i := 0; i < n-1; i++ {
			if arr[i] > arr[i+1] {
				res++
				i++
			}
		}
		fmt.Fprintln(out, res)
	}
}
