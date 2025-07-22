package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

	d := arr[1] - arr[0]
	for i := 2; i < n; i++ {
		diff := arr[i] - arr[i-1]
		d = gcd(d, diff)
	}

	var add int64
	for i := 1; i < n; i++ {
		diff := arr[i] - arr[i-1]
		add += diff/d - 1
	}
	fmt.Fprintln(out, add)
}
