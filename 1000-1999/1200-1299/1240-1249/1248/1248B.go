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

	var n int
	fmt.Fscan(in, &n)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	var x, y int64
	for i := 0; i < n/2; i++ {
		x += arr[i]
	}
	for i := n / 2; i < n; i++ {
		y += arr[i]
	}
	res := x*x + y*y
	fmt.Fprintln(out, res)
}
