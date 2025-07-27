package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
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
		var k int64
		fmt.Fscan(in, &n, &k)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		d := int64(0)
		for i := 1; i < n; i++ {
			diff := arr[i] - arr[0]
			if diff < 0 {
				diff = -diff
			}
			d = gcd(d, diff)
		}
		diffK := k - arr[0]
		if diffK%d == 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
