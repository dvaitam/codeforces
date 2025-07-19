package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	sum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		sum[i] = sum[i-1] + a[i]
	}
	ans := 0
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if sum[i]-sum[j] > (i-j)*100 {
				if i-j > ans {
					ans = i - j
				}
			}
		}
	}
	fmt.Print(ans)
}
