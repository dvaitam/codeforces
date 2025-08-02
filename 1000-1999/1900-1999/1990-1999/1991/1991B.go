package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 1000003

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n+1)
		for i := 1; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b[n] = a[n-1]
		b[1] = a[1]
		for i := 2; i < n; i++ {
			b[i] = a[i-1] | a[i]
		}
		for i := 1; i < n; i++ {
			if b[i]&b[i+1] != a[i] {
				fmt.Println(-1)
				goto end
			}
		}
		for i := 1; i <= n; i++ {
			fmt.Print(b[i], " ")
		}
		fmt.Println()
	end:
	}
}
