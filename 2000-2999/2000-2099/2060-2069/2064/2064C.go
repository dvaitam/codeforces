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
		a := make([]int64, n+2)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		posPrefix := make([]int64, n+2)
		for i := 1; i <= n; i++ {
			posPrefix[i] = posPrefix[i-1]
			if a[i] > 0 {
				posPrefix[i] += a[i]
			}
		}
		negSuffix := make([]int64, n+2)
		for i := n; i >= 1; i-- {
			negSuffix[i] = negSuffix[i+1]
			if a[i] < 0 {
				negSuffix[i] += -a[i]
			}
		}
		ans := negSuffix[1]
		for j := 1; j <= n; j++ {
			if a[j] > 0 {
				temp := posPrefix[j] + negSuffix[j+1]
				if temp > ans {
					ans = temp
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
