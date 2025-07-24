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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Ints(arr)
	for i := 0; i+2 < n; i++ {
		if arr[i]+arr[i+1] > arr[i+2] {
			fmt.Fprintln(out, "YES")
			return
		}
	}
	fmt.Fprintln(out, "NO")
}
