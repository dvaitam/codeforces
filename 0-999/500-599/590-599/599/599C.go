package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	prefixMax := make([]int, n)
	prefixMax[0] = a[0]
	for i := 1; i < n; i++ {
		if a[i] > prefixMax[i-1] {
			prefixMax[i] = a[i]
		} else {
			prefixMax[i] = prefixMax[i-1]
		}
	}

	suffixMin := make([]int, n)
	suffixMin[n-1] = a[n-1]
	for i := n - 2; i >= 0; i-- {
		if a[i] < suffixMin[i+1] {
			suffixMin[i] = a[i]
		} else {
			suffixMin[i] = suffixMin[i+1]
		}
	}

	ans := 1
	for i := 0; i < n-1; i++ {
		if prefixMax[i] <= suffixMin[i+1] {
			ans++
		}
	}
	fmt.Println(ans)
}
