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
		arr := make([]int64, n)
		var g int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if i == 0 {
				g = arr[i]
			} else {
				g = gcd(g, arr[i])
			}
		}

		var res int64
		for i := 0; i < n; i++ {
			if arr[i] != g {
				res += arr[i] - g
			}
		}
		fmt.Fprintln(out, res)
	}
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}
