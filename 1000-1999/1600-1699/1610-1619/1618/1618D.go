package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	br := bufio.NewReader(os.Stdin)
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	var t int
	fmt.Fscan(br, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(br, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(br, &arr[i])
		}
		sort.Ints(arr)
		res := 0
		for i := 0; i < k; i++ {
			a := arr[n-2*k+i]
			b := arr[n-k+i]
			res += a / b
		}
		for i := 0; i < n-2*k; i++ {
			res += arr[i]
		}
		fmt.Fprintln(bw, res)
	}
}
