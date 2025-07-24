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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Ints(arr)

	if k == 0 {
		if arr[0] == 1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, arr[0]-1)
		}
		return
	}
	if k == n {
		fmt.Fprintln(out, arr[n-1])
		return
	}
	x := arr[k-1]
	if arr[k] == x {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, x)
	}
}
