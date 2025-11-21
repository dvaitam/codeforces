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
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			sum += arr[i]
		}
		ans := sum - 2*arr[n-2]
		fmt.Fprintln(out, ans)
	}
}
