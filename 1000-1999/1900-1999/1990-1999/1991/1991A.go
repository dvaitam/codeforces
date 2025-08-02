package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var queries int
	fmt.Fscan(in, &queries)
	for i := 0; i < queries; i++ {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[j])
		}
		mx := int64(0)
		for j := 0; j < n; j += 2 {
			mx = max(mx, a[j])
		}
		fmt.Println(mx)
	}
}
