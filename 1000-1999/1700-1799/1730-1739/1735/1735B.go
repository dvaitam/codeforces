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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		minVal := arr[0]
		denom := 2*minVal - 1
		var ans int64
		for i := 1; i < n; i++ {
			ans += (arr[i] - 1) / denom
		}
		fmt.Fprintln(out, ans)
	}
}
