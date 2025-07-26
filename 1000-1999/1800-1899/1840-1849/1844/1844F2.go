package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var c int
		fmt.Fscan(in, &n, &c)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		if c >= 0 {
			sort.Ints(arr)
		} else {
			sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, arr[i])
		}
		fmt.Fprintln(out)
	}
}
